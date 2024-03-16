import {App} from './App';
import * as ServiceWorker from './ServiceWorker';
import {AppearanceProvider} from "./provider/AppearanceProvider";
import {StoreProvider} from "./provider/StoreProvider";
import {CssBaseline} from "@mui/material";
import {createRoot} from "react-dom/client";
import {SnackbarProvider} from "notistack";
import scroll from "./style/scroll.module.css"
import {StrictMode} from "react";
import {AuthProvider} from "./provider/AuthProvider";
import {LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import utc from "dayjs/plugin/utc";
import dayjs from "dayjs";

// reserve place for scroll to avoid resizing
document.documentElement.classList.add(scroll.hidden)
document.body.classList.add(scroll.show)

// extend dayjs with UTC plugin
dayjs.extend(utc)

// render react app
const container = document.getElementById("root");
const root = createRoot(container!);
root.render(
    <StrictMode>
        <LocalizationProvider dateAdapter={AdapterDayjs}>
            <AppearanceProvider>
                <AuthProvider>
                    <StoreProvider>
                        <SnackbarProvider maxSnack={3}>
                            <CssBaseline enableColorScheme/>
                            <App/>
                        </SnackbarProvider>
                    </StoreProvider>
                </AuthProvider>
            </AppearanceProvider>
        </LocalizationProvider>
    </StrictMode>
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
ServiceWorker.unregister();
