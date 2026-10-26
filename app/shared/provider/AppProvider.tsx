import {
    createTheme,
    CssBaseline,
    PaletteMode,
    Theme,
    ThemeProvider as MuiThemeProvider,
    useMediaQuery
} from "@mui/material"
import {focusManager, QueryClient, QueryClientProvider} from "@tanstack/react-query"
import {ReactQueryDevtools} from "@tanstack/react-query-devtools"
import {createContext, ReactNode, useContext, useEffect} from "react"

import {ThemeComponents} from "../helper/HelperTheme"
import {useLocalStorageState} from "../hook/LocalStorage"

export enum Mode {
    DARK = "dark",
    LIGHT = "light",
    SYSTEM = "system",
}

interface AppearanceStateType {
    mode: Mode,
    refetchOnWindowsFocus: boolean,
}

interface ThemeContextType {
    state: AppearanceStateType,
    theme: "dark" | "light",
    setTheme: (m: Mode) => void,
    toggleRefetchOnWindowsRefocus: () => void,
    info?: Theme,
}

const ThemeInitialState: AppearanceStateType = {mode: Mode.SYSTEM, refetchOnWindowsFocus: false}
const ThemeContext = createContext<ThemeContextType>({
    state: ThemeInitialState,
    theme: "light",
    toggleRefetchOnWindowsRefocus: () => void 0,
    setTheme: () => void 0,
})
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

export const MainQueryClient = new QueryClient()

export function useSettings() {
    return useContext(ThemeContext)
}

export function AppProvider(props: { children: ReactNode }) {
    const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)")
    const [state, setState] = useLocalStorageState("appearance", ThemeInitialState, true)

    const theme = getTheme(state.mode)
    const muiTheme = createTheme({palette: {mode: theme}, components: ThemeComponents})

    useEffect(handleEffectClient, [state.refetchOnWindowsFocus])
    useEffect(handleEffectThemeColor, [theme, muiTheme.palette.background.default])

    return (
        <ThemeContext value={{state, theme, setTheme, toggleRefetchOnWindowsRefocus, info: muiTheme}}>
            <QueryClientProvider client={MainQueryClient}>
                <ReactQueryDevtools/>
                <MuiThemeProvider theme={muiTheme}>
                    <CssBaseline enableColorScheme/>
                    {props.children}
                </MuiThemeProvider>
            </QueryClientProvider>
        </ThemeContext>
    )

    function setTheme(mode: Mode) {
        return setState(prevState => ({...prevState, mode}))
    }

    function toggleRefetchOnWindowsRefocus() {
        return setState(prevState => ({...prevState, refetchOnWindowsFocus: !state.refetchOnWindowsFocus}))
    }

    function handleEffectClient() {
        const enabled = state.refetchOnWindowsFocus ? "always" : false
        // TODO this works only in initial initialization after page rerender...
        //  it should be set for each query separately to make it work, probably we need to have some wrapper
        //  of useQuery or reload page?
        MainQueryClient.setDefaultOptions({queries: {refetchOnWindowFocus: enabled}})
    }

    // NOTE: index.html follows the system before the app draws, a theme picked by hand must win over it, so one tag
    //  without a media query is left; Chrome on Android paints its toolbar with it. An installed app's navigation bar
    //  takes the manifest's colour instead, one per file and read on install, so the manifest follows the theme too.
    //  The manifest url stays relative, so it resolves against `<base href>` behind `IVORY_URL_PATH`
    function handleEffectThemeColor() {
        document.head.querySelector<HTMLLinkElement>("link[rel='manifest']")
            ?.setAttribute("href", theme === "dark" ? "./manifest-dark.json" : "./manifest-light.json")
        const tags = document.head.querySelectorAll<HTMLMetaElement>("meta[name='theme-color']")
        tags.forEach((tag, index) => index === 0 ? tag.removeAttribute("media") : tag.remove())
        const tag = tags[0] ?? document.head.appendChild(Object.assign(document.createElement("meta"), {name: "theme-color"}))
        tag.content = muiTheme.palette.background.default
    }

    function getTheme(m: Mode): PaletteMode {
        switch (m) {
            case Mode.DARK:
                return "dark"
            case Mode.LIGHT:
                return "light"
            default:
                return prefersDarkMode ? "dark" : "light"
        }
    }
}
