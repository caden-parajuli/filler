/**
 * @returns boolean - Whether localStorage is available
 */
export function detectLocalStorage() {
    try {
        const x = "__storage_test__";
        localStorage.setItem(x, x);
        localStorage.removeItem(x);
        return true;
    } catch (e) {
        return (
            e instanceof DOMException &&
            e.name === "QuotaExceededError" &&
            localStorage &&
            localStorage.length !== 0
        );
    }
}
