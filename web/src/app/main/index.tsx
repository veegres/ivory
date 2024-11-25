import utc from "dayjs/plugin/utc";
import dayjs from "dayjs";
import * as ServiceWorker from "../../ServiceWorker";
import {App} from './App';
import {CssBaseline} from "@mui/material";
import {createRoot} from "react-dom/client";
import {StrictMode} from "react";
import {LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {ErrorBoundary} from "react-error-boundary";
import {SettingsProvider} from "../../provider/SettingsProvider";
import {StoreProvider} from "../../provider/StoreProvider";
import {AuthProvider} from "../../provider/AuthProvider";
import {SnackbarProvide} from "../../provider/SnackbarProvider";
import {PageErrorBox} from "../../component/view/box/PageErrorBox";
import scroll from "../../style/scroll.module.css"

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
        <SettingsProvider>
            <CssBaseline enableColorScheme/>
            <StoreProvider>
                <LocalizationProvider dateAdapter={AdapterDayjs}>
                    <AuthProvider>
                        <SnackbarProvide>
                            <ErrorBoundary fallbackRender={(e) => (<App><PageErrorBox error={e.error}/></App>)}>
                                <App />
                            </ErrorBoundary>
                        </SnackbarProvide>
                    </AuthProvider>
                </LocalizationProvider>
            </StoreProvider>
        </SettingsProvider>
    </StrictMode>
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
ServiceWorker.unregister();
