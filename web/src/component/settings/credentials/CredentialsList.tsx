import {useQuery} from "@tanstack/react-query";
import {passwordApi} from "../../../app/api";
import {ErrorAlert} from "../../view/ErrorAlert";
import {Password} from "../../../type/password";
import {InfoAlert} from "../../view/InfoAlert";
import {Collapse} from "@mui/material";
import {TransitionGroup} from "react-transition-group";
import {CredentialsItem} from "./CredentialsItem";
import {LinearProgressStateful} from "../../view/LinearProgressStateful";
import {MenuWrapperScroll} from "../menu/MenuWrapperScroll";

export function CredentialsList() {
    const query = useQuery(["credentials"], () => passwordApi.list())
    const {data: credentials, isError, error, isFetching} = query

    return (
        <>
            <LinearProgressStateful color={"inherit"} isFetching={isFetching} line/>
            {renderList()}
        </>
    )

    function renderList() {
        if (isError) return <ErrorAlert error={error}/>

        const list = Object.entries<Password>(credentials ?? {})
        if (list.length === 0) return <InfoAlert text={"There is no credentials yet"}/>

        return (
            <MenuWrapperScroll>
                <TransitionGroup appear={false}>
                    {list.map(([key, credential]) => (
                        <Collapse key={key}>
                            <CredentialsItem uuid={key} credential={credential}/>
                        </Collapse>
                    ))}
                </TransitionGroup>
            </MenuWrapperScroll>
        )
    }
}
