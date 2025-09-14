import utc from "dayjs/plugin/utc";
import dayjs from "dayjs";
import * as ServiceWorker from "../../ServiceWorker";
import {App} from './App';
import {createRoot} from "react-dom/client";
import {StrictMode} from "react";
import {LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {ErrorBoundary} from "react-error-boundary";
import {AppProvider} from "../../provider/AppProvider";
import {AuthProvider} from "../../provider/AuthProvider";
import {SnackbarProvide} from "../../provider/SnackbarProvider";
import {PageErrorBox} from "../../component/view/box/PageErrorBox";
import scroll from "../../style/scroll.module.css"
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {ReactQueryDevtools} from "@tanstack/react-query-devtools";

// extend dayjs with UTC plugin
dayjs.extend(utc)

// reserve place for scroll to avoid resizing
document.body.classList.add(scroll.hidden)
document.getElementById("root")!.classList.add(scroll.show)

export const MainQueryClient = new QueryClient()

// render react app
const container = document.getElementById("root")
const root = createRoot(container!)
root.render(
    <StrictMode>
        <QueryClientProvider client={MainQueryClient}>
            <ReactQueryDevtools/>
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
        </QueryClientProvider>
    </StrictMode>
)

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
ServiceWorker.unregister();
