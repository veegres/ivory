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

// always show scroll bar to avoid resizing
document.body.classList.add(scroll.show)

// render react app
const container = document.getElementById("root");
const root = createRoot(container!);
root.render(
    <StrictMode>
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
    </StrictMode>
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
ServiceWorker.unregister();
