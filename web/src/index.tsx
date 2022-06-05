import React from 'react';
import ReactDOM from 'react-dom';
import {App} from './App';
import * as ServiceWorker from './ServiceWorker';
import {QueryClient, QueryClientProvider} from "react-query";
import {ReactQueryDevtools} from 'react-query/devtools'
import {ThemeProvider} from "./provider/ThemeProvider";
import {StoreProvider} from "./provider/StoreProvider";
import {CssBaseline} from "@mui/material";

ReactDOM.render(
    <React.StrictMode>
        <QueryClientProvider client={new QueryClient()}>
            <ThemeProvider>
                <StoreProvider>
                    <CssBaseline enableColorScheme />
                    <App/>
                </StoreProvider>
            </ThemeProvider>
            <ReactQueryDevtools/>
        </QueryClientProvider>
    </React.StrictMode>,
    document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
ServiceWorker.unregister();
