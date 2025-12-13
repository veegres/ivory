import {useSnackbar} from "../provider/SnackbarProvider"

export function useQueryParamErrorHandler() {
    const snackbar = useSnackbar()
    const error = new URLSearchParams(window.location.search).get("error")
    if (error) {
        history.replaceState(null, "",  window.location.origin + window.location.pathname)
        // NOTE: we need to wait for the snackbar to be rendered, that is why we use setTimeout
        setTimeout(() => snackbar(`EXTERNAL - ERROR, ${error}`, "error"), 100)
    }
}