import {useMutation, useQuery} from "react-query";
import {credentialApi} from "../../app/api";
import {Error} from "../view/Error";
import {AxiosError} from "axios";
import {Credential, Style} from "../../app/types";
import {Info} from "../view/Info";
import {
    Avatar,
    Box,
    CircularProgress, Collapse,
    IconButton,
    TextField,
    Tooltip
} from "@mui/material";
import React from "react";
import {Add, Delete, Edit, Save, Storage} from "@mui/icons-material";
import {shortUuid} from "../../app/utils";
import {LinearProgressStateful} from "../view/LinearProgressStateful";
import {TransitionGroup} from "react-transition-group";

const SX = {
    field: { margin: "0 5px", width: "100%" },
    uuid: { fontSize: "8px" },
    progress: { margin: "20px 0" },
    item: { display: "flex", alignItems: "center", justifyContent: "space-between", gap: "15px", margin: "3px 0" },
    cred: { display: "flex", gap: "15px" }
}

const style: Style = {
    transition: {}
}

export function CredentialsContent() {
    const { data: credentials, isError, error, isFetching, refetch } = useQuery("credentials", credentialApi.get)
    const deleteCredentials = useMutation(credentialApi.delete, { onSuccess: refetch })

    if (isError) return <Error error={error as AxiosError}/>

    return (
        <Box>
            {renderAdd()}
            <LinearProgressStateful sx={SX.progress} color={"inherit"} isFetching={isFetching} line />
            {renderList()}
        </Box>
    )

    function renderList() {
        const list = Object.entries<Credential>(credentials ?? {})
        if (list.length === 0) return <Info text={"There is no credentials yet"}/>

        return (
            <TransitionGroup style={style.transition}>
                {list.map(([key, value]) => (
                    <Collapse key={key}>{renderItem(key, value)}</Collapse>
                ))}
            </TransitionGroup>
        )
    }

    function renderAdd() {
        return (
            <Box sx={SX.item}>
                <Tooltip title={"postgres"} placement={"right"}>
                    <Avatar><Add/></Avatar>
                </Tooltip>
                <Box sx={SX.cred}>
                    <TextField size={"small"} variant={"outlined"} label={`username`} />
                    <TextField size={"small"} variant={"outlined"} label={"password"} type={"password"} />
                </Box>
                <Box>
                    <IconButton size={"small"} onClick={() => {}}>
                        <Save />
                    </IconButton>
                </Box>
            </Box>
        )
    }

    function renderItem(uuid: string, credential: Credential) {
        const { username, password } = credential
        const name = username === "" ? "empty" : username
        const isLoading = deleteCredentials.isLoading && deleteCredentials.variables === uuid
        return (
            <Box sx={SX.item}>
                <Tooltip title={"postgres"} placement={"right"}>
                    <Avatar><Storage/></Avatar>
                </Tooltip>
                <Box sx={SX.cred}>
                    <TextField size={"small"} variant={"standard"} label={`ID: ${shortUuid(uuid)}`} value={name} disabled/>
                    <TextField size={"small"} variant={"standard"} label={"password"} value={password} disabled/>
                </Box>
                <Box>
                    <IconButton size={"small"} onClick={() => {}}>
                        <Edit />
                    </IconButton>
                    <IconButton size={"small"} onClick={() => deleteCredentials.mutate(uuid)} disabled={isLoading}>
                        {!isLoading ? <Delete/> : <CircularProgress size={25}/>}
                    </IconButton>
                </Box>
            </Box>
        )
    }
}
