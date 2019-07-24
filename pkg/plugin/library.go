package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/jkcfg/jk/pkg/plugin/renderer"
	"github.com/pkg/errors"
)

// LibraryOptions are input parameters used by the library constructor.
type LibraryOptions struct {
	// Verbose indicates if the library should print out what it is doing.
	Verbose bool
	// PluginDirectory is where the downloaded plugins are stored.
	PluginDirectory string
}

type phase int

const (
	new phase = iota
	starting
	running
	failure
)

type state struct {
	phase     phase
	phaseMu   sync.Mutex
	phaseCond *sync.Cond
	phaseErr  error

	proto plugin.ClientProtocol
}

func newState(p phase) *state {
	s := &state{
		phase: p,
	}
	s.phaseCond = sync.NewCond(&s.phaseMu)
	return s
}

func (s *state) waitForRunning() error {
	s.phaseMu.Lock()
	defer s.phaseMu.Unlock()

	for {
		if s.phase == running {
			return nil
		}
		if s.phase == failure {
			return s.phaseErr
		}
		s.phaseCond.Wait()
	}
}

func (s *state) setPhase(p phase) {
	s.phaseMu.Lock()
	s.phase = p
	s.phaseMu.Unlock()
	s.phaseCond.Broadcast()
}

func (s *state) setError(err error) {
	s.phaseMu.Lock()
	s.phaseErr = err
	s.phaseMu.Unlock()
	s.phaseCond.Broadcast()
}

// Library is a library of plugins. This is a factory object handling the full
// life cycle of plugins: download, creation and termination of plugin
// processes.
type Library struct {
	LibraryOptions

	lock    sync.Mutex // protects the plugins map.
	plugins map[string]*state
}

// NewLibrary creates a new plugin library.
func NewLibrary(opts LibraryOptions) *Library {
	return &Library{
		LibraryOptions: opts,
		plugins:        make(map[string]*state),
	}
}

// pluginMap is the map of plugins we can dispense.
var rendererPluginMap = map[string]plugin.Plugin{
	"renderer": &renderer.Plugin{},
}

func fetchLocalInfo(path string) (*Info, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var info Info
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, errors.Wrapf(err, "parsing %q", path)
	}
	return &info, nil
}

func fetchInfo(pluginInfo string) (*Info, error) {
	url, err := url.Parse(pluginInfo)
	if err != nil {
		return nil, err
	}

	switch url.Scheme {
	case "":
		return fetchLocalInfo(url.Path)
	default:
		return nil, fmt.Errorf("unknown scheme %q", url.Scheme)
	}
}

// get retrieves a plugin from the library, returning an object that can be used
// to issue remote procedure calls.
func (l *Library) get(kind string, pluginInfo string) (plugin.ClientProtocol, error) {
	var phase phase

	l.lock.Lock()
	s := l.plugins[pluginInfo]
	if s == nil {
		// phase will be:
		//  - 'new' for first get() invocation (first call for each distinct
		//     pluginInfo).
		//  - 'starting' for get() calls that arrive between this point and the
		//     end of the first call.
		s = newState(starting)
		l.plugins[pluginInfo] = s
		phase = new
	} else {
		phase = s.phase
	}
	l.lock.Unlock()

	// Plugin is running, just return the client object.
	if phase == running {
		return s.proto, nil
	}

	// Plugin is stuck in failure state, return the associated error.
	if phase == failure {
		return nil, s.phaseErr
	}

	// Plugin is still starting wait for it to be running.
	if phase == starting {
		err := s.waitForRunning()
		return s.proto, err
	}

	// 'new' phase, starting the plugin.

	// When returning from get, whether the plugin has been successfully started
	// or not, we need to unlock everyone that is waiting for the running state.
	var err error
	defer func() {
		if err == nil {
			s.setPhase(running)
		} else {
			s.setError(err)
		}
	}()

	info, err := fetchInfo(pluginInfo)
	if err != nil {
		return nil, err
	}

	binary := info.binary()
	if binary == "" {
		err = errors.New("no plugin binary for your os/processor")
		return nil, err
	}

	var proto plugin.ClientProtocol

	switch kind {
	case "renderer":
		if l.Verbose {
			fmt.Printf("starting plugin %q\n", binary)
		}

		// Start by launching the plugin process.
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: renderer.RendererV1,
			Managed:         true,
			Plugins:         rendererPluginMap,
			Logger: hclog.New(&hclog.LoggerOptions{
				Level:  hclog.Info,
				Output: os.Stderr,
			}),
			Cmd: exec.Command(binary),
		})

		proto, err = client.Client()
	default:
		err = fmt.Errorf("unknown kind %q", kind)
	}

	if err != nil {
		return nil, err
	}

	l.lock.Lock()
	s = l.plugins[pluginInfo]
	s.proto = proto
	s.phase = running
	l.lock.Unlock()

	return proto, nil
}

// GetRenderer materializes a plugin URI pointing at JSON plugin.Info into a
// callable interface.
func (l *Library) GetRenderer(pluginInfo string) (renderer.Renderer, error) {
	client, err := l.get("renderer", pluginInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching %s", pluginInfo)
	}

	raw, err := client.Dispense("renderer")
	if err != nil {
		return nil, err
	}
	return raw.(renderer.Renderer), nil
}

// Close terminates the library, terminating all plugins.
func (l *Library) Close() {
	plugin.CleanupClients()
	l.plugins = make(map[string]*state)
}
