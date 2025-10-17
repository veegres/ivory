import {AuthType} from "../config/type";

export interface AuthInfo {
    type: AuthType,
    authorised: boolean,
    error: string,
}