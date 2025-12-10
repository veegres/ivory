import {LocalizationProvider} from "@mui/x-date-pickers"
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs"
import dayjs from "dayjs"
import utc from "dayjs/plugin/utc"
import {StrictMode} from "react"
import {createRoot} from "react-dom/client"
import {ErrorBoundary} from "react-error-boundary"

import {PageErrorBox} from "../../component/view/box/PageErrorBox"
import {AppProvider} from "../../provider/AppProvider"
import {AuthProvider} from "../../provider/AuthProvider"
import {SnackbarProvide} from "../../provider/SnackbarProvider"
import scroll from "../../style/scroll.module.css"
import {App} from "./App"

// extend dayjs with UTC plugin
dayjs.extend(utc)

// reserve place for scroll to avoid resizing
document.body.classList.add(scroll.hidden)
document.getElementById("root")!.classList.add(scroll.show)

// render react app
const container = document.getElementById("root")
const root = createRoot(container!)
root.render(
    <StrictMode>
        <AppProvider>
            <LocalizationProvider dateAdapter={AdapterDayjs}>
                <AuthProvider>
                    <SnackbarProvide>
                        <ErrorBoundary fallbackRender={(e) => (<App><PageErrorBox error={e.error}/></App>)}>
                            <App/>
                        </ErrorBoundary>
                    </SnackbarProvide>
                </AuthProvider>
            </LocalizationProvider>
        </AppProvider>
    </StrictMode>
)

