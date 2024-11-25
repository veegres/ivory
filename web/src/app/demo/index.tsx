import {CssBaseline} from "@mui/material";
import {createRoot} from "react-dom/client";
import scroll from "../../style/scroll.module.css"
import {StrictMode} from "react";
import {LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import utc from "dayjs/plugin/utc";
import dayjs from "dayjs";
import {App} from "./App";
import {SettingsProvider} from "../../provider/SettingsProvider";

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
            <LocalizationProvider dateAdapter={AdapterDayjs}>
                <App />
            </LocalizationProvider>
        </SettingsProvider>
    </StrictMode>
)
