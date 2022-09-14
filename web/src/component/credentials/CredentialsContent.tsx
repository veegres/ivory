import {useQuery} from "@tanstack/react-query";
import {credentialApi} from "../../app/api";
import {ErrorAlert} from "../view/ErrorAlert";
import {AxiosError} from "axios";
import {Credential} from "../../app/types";
import {InfoAlert} from "../view/InfoAlert";
import {Box, Collapse} from "@mui/material";
import React from "react";
import {LinearProgressStateful} from "../view/LinearProgressStateful";
import {TransitionGroup} from "react-transition-group";
import {CredentialsItem} from "./CredentialsItem";
import {CredentialsNew} from "./CredentialsNew";
import scroll from "../../style/scroll.module.css"

const SX = {
    field: { margin: "0 5px", width: "100%" },
    uuid: { fontSize: "8px" },
    progress: { margin: "5px 0 20px" },
    cred: { display: "flex", gap: "15px" },
    list: { maxHeight: "500px", overflowY: "auto" },
}

export function CredentialsContent() {
    const { data: credentials, isError, error, isFetching } = useQuery(["credentials"], () => credentialApi.get())
    if (isError) return <ErrorAlert error={error as AxiosError}/>

    return (
        <Box>
            <CredentialsNew />
            <LinearProgressStateful sx={SX.progress} color={"inherit"} isFetching={isFetching} line />
            {renderList()}
        </Box>
    )

    function renderList() {
        const list = Object.entries<Credential>(credentials ?? {})
        if (list.length === 0) return <InfoAlert text={"There is no credentials yet"}/>

        return (
            <Box sx={SX.list} className={scroll.tiny}>
                <TransitionGroup>
                    {list.map(([uuid, credential]) => (
                        <Collapse key={uuid}>
                            <CredentialsItem uuid={uuid} credential={credential}/>
                        </Collapse>
                    ))}
                </TransitionGroup>
            </Box>
        )
    }
}
