import React, {createContext, ReactNode, useContext, useState} from "react";
import {createTheme, CssBaseline, PaletteMode, Theme, ThemeProvider as MuiThemeProvider} from "@mui/material";

interface ThemeContextType {
    mode?: "dark" | "light",
    toggle: () => void
    info?: Theme,
}
const ThemeContext = createContext<ThemeContextType>({
    mode: undefined,
    toggle: () => {
    }
});

export function useTheme() {
    return useContext(ThemeContext)
}

export function ThemeProvider(props: { children: ReactNode }) {
    const [mode, setMode] = useState(getMode());
    const muiTheme = createTheme({palette: {mode}})
    return (
        <ThemeContext.Provider value={{ mode, toggle, info: muiTheme }}>
            <MuiThemeProvider theme={muiTheme}>
                <CssBaseline/>
                {props.children}
            </MuiThemeProvider>
        </ThemeContext.Provider>
    );

    function getMode() {
        const storageMode = localStorage.getItem('theme')
        return storageMode == null ? 'light' : storageMode as PaletteMode
    }

    function toggle() {
        setMode((prevMode) => {
            const currentMode = prevMode === 'light' ? 'dark' : 'light'
            localStorage.setItem('theme', currentMode)
            return currentMode
        });
    }
}
