import {createContext, ReactNode, useContext, useEffect} from "react";
import {api} from "../app/api";
import {useLocalStorageState} from "../hook/LocalStorage";
import {useQueryClient} from "@tanstack/react-query";

const AuthContext = createContext({
    setToken: (v: string) => {},
    logout: () => {},
})

export function useAuth() {
    return useContext(AuthContext)
}

export function AuthProvider(props: { children: ReactNode }) {
    const [token, setToken] = useLocalStorageState("token", "", true);
    const queryClient = useQueryClient();

    useEffect(handleEffectTokenUpdate, [token, queryClient])

    return (
        <AuthContext.Provider value={{setToken: setTokenEncrypt, logout}}>
            {props.children}
        </AuthContext.Provider>
    )

    function handleEffectTokenUpdate() {
        console.log(window.atob(token))
        if (token !== "") api.defaults.headers.common["Authorization"] = `Bearer ${window.atob(token)}`
        else delete api.defaults.headers.common["Authorization"]
        queryClient.refetchQueries(["info"]).then()
    }

    function setTokenEncrypt(v: string){
        setToken(window.btoa(v))
    }

    function logout(){
        setToken("")
    }
}
