import {createContext, ReactNode, useContext} from "react";
import {api} from "../app/api";
import {useLocalStorageState} from "../hook/LocalStorage";
import {useQueryClient} from "@tanstack/react-query";

interface AuthContextType {
    setToken: (v: string) => void,
    logout: () => void,
}
const AuthContext = createContext<AuthContextType>({
    setToken: () => void 0,
    logout: () => void 0,
})

export function useAuth() {
    return useContext(AuthContext)
}

export function AuthProvider(props: { children: ReactNode }) {
    const [token, setToken] = useLocalStorageState("token", "", true);
    const queryClient = useQueryClient();

    if (token) api.defaults.headers.common["Authorization"] = `Bearer ${window.atob(token)}`
    else delete api.defaults.headers.common["Authorization"]
    queryClient.refetchQueries({queryKey: ["info"]}).then()

    return (
        <AuthContext.Provider value={{setToken: setTokenEncrypt, logout}}>
            {props.children}
        </AuthContext.Provider>
    )

    function setTokenEncrypt(v: string){
        setToken(window.btoa(v))
    }

    function logout(){
        setToken("")
    }
}
