package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func mustParseRef(img string) name.Reference {
	ref, err := name.ParseReference(img)
	if err != nil {
		log.Fatalf("could not parse reference from %q: %w", img, err)
	}
	return ref
}

func main() {
	old := flag.String("old", "", "old base")
	nu := flag.String("new", "", "new base")
	orig := flag.String("original", "", "image to rebase")
	out := flag.String("out", "", "result tarball")

	flag.Parse()

	origImg, err := daemon.Image(mustParseRef(*orig))
	if err != nil {
		log.Fatalf("pulling %s: %v", *orig, err)
	}

	oldBaseImg, err := daemon.Image(mustParseRef(*old))
	if err != nil {
		log.Fatalf("pulling %s: %v", *old, err)
	}

	newBaseImg, err := daemon.Image(mustParseRef(*nu))
	if err != nil {
		log.Fatalf("pulling %s: %v", *nu, err)
	}

	img, err := mutate.Rebase(origImg, oldBaseImg, newBaseImg)
	if err != nil {
		log.Fatalf("rebasing: %v", err)
	}

	image := (*out)[:len(*out)-4] + ":latest"
	newTag, err := name.NewTag(image)
	if err != nil {
		log.Fatalf("could not create image ref %q: %w", image, err)
	}
	f, err := os.Create(*out)
	if err != nil {
		log.Fatalf("could not create file %q: %w", *out, err)
	}
	defer f.Close()

	if err := tarball.Write(newTag, img, f); err != nil {
		log.Fatalf("could not write image to %v: %w", *out, err)
	}

	digest, err := img.Digest()
	if err != nil {
		log.Fatalf("digesting rebased: %v", err)
	}
	fmt.Printf("Wrote image %s@%s to file %s\n", image, digest.String(), *out)
}
