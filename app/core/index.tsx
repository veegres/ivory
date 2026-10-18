import {LocalizationProvider} from "@mui/x-date-pickers"
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs"
import dayjs from "dayjs"
import utc from "dayjs/plugin/utc"
import {StrictMode} from "react"
import {createRoot} from "react-dom/client"
import {ErrorBoundary} from "react-error-boundary"

import {ManagementApi} from "../features/management/api/ManagementRouter"
import {PageErrorBox} from "../shared/component/box/PageErrorBox"
import {AppProvider, MainQueryClient} from "../shared/provider/AppProvider"
import {AuthProvider} from "../shared/provider/AuthProvider"
import {SnackbarProvide} from "../shared/provider/SnackbarProvider"
import background from "../shared/style/background.module.css"
import scroll from "../shared/style/scroll.module.css"
import {App} from "./App"

// extend dayjs with UTC plugin
dayjs.extend(utc)

// reserve place for scroll to avoid resizing
document.body.classList.add(scroll.hidden)
document.getElementById("root")!.classList.add(scroll.show)
document.getElementById("root")!.classList.add(background.inherit)

// render react app
const container = document.getElementById("root")
const root = createRoot(container!)
root.render(
    <StrictMode>
        <AppProvider>
            <LocalizationProvider dateAdapter={AdapterDayjs}>
                <AuthProvider>
                    <SnackbarProvide>
                        <ErrorBoundary fallbackRender={(e) => <PageErrorBox error={e.error} description={getCrashDescription()}/>}>
                            <App/>
                        </ErrorBoundary>
                    </SnackbarProvide>
                </AuthProvider>
            </LocalizationProvider>
        </AppProvider>
    </StrictMode>
)

// NOTE: App reads the info answer on every render, so one that is not Ivory's is the cause worth naming
function getCrashDescription() {
    const info = MainQueryClient.getQueryData<Record<string, unknown>>(ManagementApi.info.key())
    if (info && (!info.auth || !info.config || !info.secret || !info.version)) {
        return "The server behind /api answered, but not with Ivory's application info. " +
            "Another service is probably listening on the port the Ivory server should use, " +
            "or the server is a different Ivory version than this page."
    }
    return "Ivory failed to render this page. Reloading usually helps."
}
