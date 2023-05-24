import {createContext, ReactNode, useContext} from "react";
import {api} from "../app/api";

const AuthContext = createContext({
    setToken: (v: string) => {},
    logout: () => {},
})

export function useAuth() {
    return useContext(AuthContext)
}

export function AuthProvider(props: { children: ReactNode }) {
    const token = localStorage.getItem("token")
    if (token) setToken(token)
    return (
        <AuthContext.Provider value={{setToken, logout}}>{props.children}</AuthContext.Provider>
    )

    function setToken(v: string){
        localStorage.setItem("token", v)
        api.defaults.headers.common["Authorization"] = `Bearer ${v}`
    }

    function logout(){
        localStorage.removeItem("token")
        delete api.defaults.headers.common["Authorization"]
    }
}
