import {createContext, ReactNode, useContext, useEffect} from "react";
import {
    createTheme,
    CssBaseline,
    PaletteMode,
    Theme,
    ThemeProvider as MuiThemeProvider,
    useMediaQuery
} from "@mui/material";
import {useLocalStorageState} from "../hook/LocalStorage";
import {focusManager, QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {ReactQueryDevtools} from "@tanstack/react-query-devtools";

export enum Mode {
    DARK = "dark",
    LIGHT = "light",
    SYSTEM = "system",
}

interface AppearanceStateType {
    mode: Mode,
    refetchOnWindowsFocus: boolean,
    uncheckInstanceBlockOnClusterChange: boolean,
}

interface ThemeContextType {
    state: AppearanceStateType,
    theme: "dark" | "light",
    setTheme: (m: Mode) => void,
    toggleRefetchOnWindowsRefocus: () => void,
    toggleUncheckInstanceBlockOnClusterChange: () => void,
    info?: Theme,
}

const ThemeInitialState: AppearanceStateType = {mode: Mode.SYSTEM, refetchOnWindowsFocus: false, uncheckInstanceBlockOnClusterChange: true}
const ThemeContext = createContext<ThemeContextType>({
    state: ThemeInitialState,
    theme: "light",
    toggleRefetchOnWindowsRefocus: () => void 0,
    toggleUncheckInstanceBlockOnClusterChange: () => void 0,
    setTheme: () => void 0,
})
const client = new QueryClient()
focusManager.setEventListener(handleFocus => {
    if (typeof window !== "undefined" && window.addEventListener) {
        window.addEventListener("focus", () => handleFocus(), false)
        window.addEventListener("visibilitychange", () => handleFocus(), false)
        return () => {
            window.removeEventListener("focus", () => handleFocus())
            window.removeEventListener("visibilitychange", () => handleFocus())
        }
    }
})

export function useSettings() {
    return useContext(ThemeContext)
}

export function SettingsProvider(props: { children: ReactNode }) {
    const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)')
    const [state, setState] = useLocalStorageState("appearance", ThemeInitialState, true)

    const theme = getTheme(state.mode)
    const muiTheme = createTheme({palette: {mode: theme}})

    useEffect(handleEffectClient, [state.refetchOnWindowsFocus])

    return (
        <ThemeContext value={{state, theme, setTheme, toggleRefetchOnWindowsRefocus, toggleUncheckInstanceBlockOnClusterChange: toggleUncheckInstanceBlockOnClusterSelect, info: muiTheme}}>
            <QueryClientProvider client={client}>
                <MuiThemeProvider theme={muiTheme}>
                    <CssBaseline enableColorScheme/>
                    {props.children}
                </MuiThemeProvider>
                <ReactQueryDevtools/>
            </QueryClientProvider>
        </ThemeContext>
    );

    function setTheme(mode: Mode) {
        return setState(prevState => ({...prevState, mode}))
    }

    function toggleRefetchOnWindowsRefocus() {
        return setState(prevState => ({...prevState, refetchOnWindowsFocus: !state.refetchOnWindowsFocus}))
    }

    function toggleUncheckInstanceBlockOnClusterSelect() {
        return setState(prevState => ({...prevState, uncheckInstanceBlockOnClusterChange: !state.uncheckInstanceBlockOnClusterChange}))
    }

    function handleEffectClient() {
        client.setDefaultOptions({queries: {refetchOnWindowFocus: state.refetchOnWindowsFocus ? "always" : false}})
    }

    function getTheme(m: Mode): PaletteMode {
        switch (m) {
            case Mode.DARK: return "dark"
            case Mode.LIGHT: return "light"
            default: return prefersDarkMode ? "dark" : "light"
        }
    }
}
