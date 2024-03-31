import {App} from './App';
import * as ServiceWorker from './ServiceWorker';
import {AppearanceProvider} from "./provider/AppearanceProvider";
import {StoreProvider} from "./provider/StoreProvider";
import {CssBaseline} from "@mui/material";
import {createRoot} from "react-dom/client";
import scroll from "./style/scroll.module.css"
import {StrictMode} from "react";
import {AuthProvider} from "./provider/AuthProvider";
import {LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import utc from "dayjs/plugin/utc";
import dayjs from "dayjs";
import {ErrorBoundary} from "react-error-boundary";
import {PageErrorBox} from "./component/view/box/PageErrorBox";
import {SnackbarProvide} from "./provider/SnackbarProvider";

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
        <AppearanceProvider>
            <ErrorBoundary fallbackRender={(e) => (<PageErrorBox error={e.error}/>)}>
                <LocalizationProvider dateAdapter={AdapterDayjs}>
                    <AuthProvider>
                        <StoreProvider>
                            <SnackbarProvide>
                                <CssBaseline enableColorScheme/>
                                <App/>
                            </SnackbarProvide>
                        </StoreProvider>
                    </AuthProvider>
                </LocalizationProvider>
            </ErrorBoundary>
        </AppearanceProvider>
    </StrictMode>
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
ServiceWorker.unregister();
