import {CredentialsNew} from "./CredentialsNew";
import {MenuWrapper} from "../menu/MenuWrapper";
import {CredentialsList} from "./CredentialsList";
import {useQuery} from "@tanstack/react-query";
import {PasswordApi} from "../../../app/api";
import {LinearProgressStateful} from "../../view/progress/LinearProgressStateful";

export function Credentials() {
    const query = useQuery({queryKey: ["credentials"], queryFn: () => PasswordApi.list()})
    const {data, error, isFetching} = query

    return (
        <MenuWrapper>
            <CredentialsNew/>
            <LinearProgressStateful color={"inherit"} isFetching={isFetching} line/>
            <CredentialsList credentials={data} error={error}/>
        </MenuWrapper>
    )
}
