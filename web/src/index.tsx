import {App} from './App';
import * as ServiceWorker from './ServiceWorker';
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {ReactQueryDevtools} from "@tanstack/react-query-devtools"
import {ThemeProvider} from "./provider/ThemeProvider";
import {StoreProvider} from "./provider/StoreProvider";
import {CssBaseline} from "@mui/material";
import {createRoot} from "react-dom/client";
import {SnackbarProvider} from "notistack";
import scroll from "./style/scroll.module.css"
import {StrictMode} from "react";

// always show scroll bar to avoid resizing
document.body.classList.add(scroll.show)

// render react app
const container = document.getElementById("root");
const root = createRoot(container!);
root.render(
    <StrictMode>
        <QueryClientProvider client={new QueryClient()}>
            <ThemeProvider>
                <StoreProvider>
                    <SnackbarProvider maxSnack={3}>
                        <CssBaseline enableColorScheme />
                        <App/>
                    </SnackbarProvider>
                </StoreProvider>
            </ThemeProvider>
            <ReactQueryDevtools/>
        </QueryClientProvider>
    </StrictMode>
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
ServiceWorker.unregister();
