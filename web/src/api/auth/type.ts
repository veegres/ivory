import {AuthType} from "../config/type";

export interface Login {
    username: string,
    password: string,
}

export interface AuthInfo {
    type: AuthType,
    authorised: boolean,
    error: string,
}