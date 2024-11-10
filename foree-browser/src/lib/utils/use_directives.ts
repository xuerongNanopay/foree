import { assertIsNode } from "./typescript";

export function clickOutside(e: HTMLElement, cb: VoidFunc) {
    function onClick(event: Event) {
        assertIsNode(event.target)
        if (!e.contains(event.target)) {
            cb();
        }
    }
    
    document.body.addEventListener('click', onClick);
    
    return {
        update(newCB: VoidFunc) {
            cb = newCB
        },
        destroy() {
            document.body.removeEventListener('click', onClick);
        }
    }
}

// export default {
//     clickOutside,
// }