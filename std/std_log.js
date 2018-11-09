import { write } from 'std_write.js';

function log(value, format) {
    write(value, "", format)
}

export {
    log
}
