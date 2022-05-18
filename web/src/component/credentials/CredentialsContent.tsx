import {useMutation, useQuery} from "react-query";
import {credentialApi} from "../../app/api";
import {Error} from "../view/Error";
import {AxiosError} from "axios";
import {Credential} from "../../app/types";
import {Info} from "../view/Info";
import {
    Avatar,
    Box,
    CircularProgress,
    IconButton, LinearProgress,
    List,
    ListItem,
    ListItemAvatar,
    ListItemText,
    TextField,
    Tooltip
} from "@mui/material";
import React from "react";
import {Delete, Lock} from "@mui/icons-material";
import {shortUuid} from "../../app/utils";

const SX = {
    field: { margin: "0 5px", width: "100%" },
    uuid: { fontSize: "8px" },
}

export function CredentialsContent() {
    const { data: credentials, isError, error, isFetching, refetch } = useQuery("credentials", credentialApi.get)
    const deleteCredentials = useMutation(credentialApi.delete, { onSuccess: refetch })

    if (isError) return <Error error={error as AxiosError}/>
    if (!credentials) return <Error error={"No data"}/>

    const list = Object.entries<Credential>(credentials)
    if (list.length === 0) return <Info text={"There is no credentials yet"}/>

    return (
        <>
            {isFetching ? <LinearProgress/> : null}
            <List disablePadding>
                {list.map(([key, value]) => (
                    <ListItem key={key} secondaryAction={renderSecondaryAction(key)}>
                        <ListItemAvatar >{renderAvatar(key)}</ListItemAvatar>
                        <ListItemText primary={renderSecondary(value)}/>
                    </ListItem>
                ))}
            </List>
        </>
    )


    function renderAvatar(uuid: string) {
        return (
            <Tooltip title={uuid} placement={"right"}>
                <Box display={"flex"} flexDirection={"column"} alignItems={"center"}>
                    <Avatar><Lock/></Avatar>
                    <Box sx={SX.uuid}>{shortUuid(uuid)}</Box>
                </Box>
            </Tooltip>
        )
    }

    function renderSecondary(credential: Credential) {
        const { username, password } = credential
        const name = username === "" ? "empty" : username
        return (
            <Box display={"flex"} justifyContent={"space-evenly"}>
                <TextField size={"small"} variant={"standard"} label={"username"} value={name} disabled/>
                <TextField size={"small"} variant={"standard"} label={"password"} value={password} disabled/>
            </Box>
        )
    }

    function renderSecondaryAction(uuid: string) {
        if (deleteCredentials.isLoading && deleteCredentials.variables === uuid) return <CircularProgress size={25}/>

        return (
            <IconButton onClick={() => deleteCredentials.mutate(uuid)}>
                <Delete/>
            </IconButton>
        )
    }
}
