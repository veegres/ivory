import React from 'react';
import ReactDOM from 'react-dom';
import { App } from './App';
import * as ServiceWorker from './ServiceWorker';
import {QueryClient, QueryClientProvider} from "react-query";

ReactDOM.render(
    <React.StrictMode>
        <QueryClientProvider client={new QueryClient()}>
            {/*<ThemeProvider theme>*/}
                <App />
            {/*</ThemeProvider>*/}
        </QueryClientProvider>
    </React.StrictMode>,
    document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
ServiceWorker.unregister();
