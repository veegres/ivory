import {createContext, ReactNode, useContext, useEffect} from "react";
import {createTheme, CssBaseline, Theme, ThemeProvider as MuiThemeProvider} from "@mui/material";
import {useLocalStorageState} from "../hook/LocalStorage";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {ReactQueryDevtools} from "@tanstack/react-query-devtools";

interface ThemeStateType {
    mode: "dark" | "light",
    refetchOnWindowsFocus: boolean,
}
interface ThemeContextType {
    state: ThemeStateType
    setTheme: (s: "dark" | "light") => void,
    toggleRefetchOnWindowsRefocus: () => void,
    info?: Theme,
}

const ThemeInitialState: ThemeStateType = { mode: "light", refetchOnWindowsFocus: false }
const ThemeContext = createContext<ThemeContextType>({
    state: ThemeInitialState,
    toggleRefetchOnWindowsRefocus: () => void 0,
    setTheme: () => void 0,
})
const client = new QueryClient()

export function useAppearance() {
    return useContext(ThemeContext)
}

export function AppearanceProvider(props: { children: ReactNode }) {
    const [state, setState] = useLocalStorageState("appearance", ThemeInitialState);
    const muiTheme = createTheme({palette: {mode: state.mode}})
    useEffect(handleEffectClient, [state.refetchOnWindowsFocus])
    return (
        <ThemeContext.Provider value={{state, setTheme, toggleRefetchOnWindowsRefocus, info: muiTheme}}>
            <QueryClientProvider client={client}>
                <MuiThemeProvider theme={muiTheme}>
                    <CssBaseline/>
                    {props.children}
                </MuiThemeProvider>
                <ReactQueryDevtools/>
            </QueryClientProvider>
        </ThemeContext.Provider>
    );

    function setTheme(m: "dark" | "light") {
        return setState(prevState => ({...prevState, mode: m}))
    }

    function toggleRefetchOnWindowsRefocus() {
        return setState(prevState => ({...prevState, refetchOnWindowsFocus: !state.refetchOnWindowsFocus}))
    }

    function handleEffectClient() {
        client.setDefaultOptions({queries: {refetchOnWindowFocus: state.refetchOnWindowsFocus ? "always" : false}})
    }
}
