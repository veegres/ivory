import {Box, Button, Grid} from "@mui/material";
import React, {ReactNode, useState} from "react";
import {useMutation, useQuery} from "@tanstack/react-query";
import {infoApi, secretApi} from "../../app/api";
import {randomUnicodeAnimal} from "../../app/utils";
import {LinearProgressStateful} from "../view/LinearProgressStateful";
import select from "../../style/select.module.css";
import {useMutationOptions} from "../../app/hooks";

const SX = {
    box: { height: "100%", width: "30%", minWidth: "500px" },
    button: { margin: "0 10px" },
    header: { fontSize: '35px', fontWeight: 900, fontFamily: 'monospace', margin: "20px 0", cursor: "pointer" },
    buttons: { margin: "8px 0" }
}

type Props = {
    children: ReactNode
    keyWord: string
    refWord: string
    clean: boolean
    header: string
}

export function Secret(props: Props) {
    const { keyWord, refWord, children, clean, header } = props
    const [animal, setAnimal] = useState(randomUnicodeAnimal())

    const info = useQuery(["info"], infoApi.get);

    const setReqOptions = useMutationOptions(["info"])
    const setReq = useMutation(secretApi.set, setReqOptions)
    const cleanReqOptions = useMutationOptions(["info"])
    const cleanReq = useMutation(secretApi.clean, cleanReqOptions)

    const fetching = cleanReq.isLoading || setReq.isLoading || info.isFetching

    return (
        <Grid container sx={SX.box} direction={"column"} alignItems={"center"} justifyContent={"center"}>
            <Box sx={SX.header} className={select.none} onClick={() => setAnimal(randomUnicodeAnimal())}>
                {header} {animal}
            </Box>
            {children}
            <LinearProgressStateful isFetching={fetching} line color={"inherit"} />
            <Box sx={SX.buttons}>
                {!clean ? null : (
                    <Button
                        sx={SX.button}
                        variant={"contained"}
                        disabled={fetching}
                        onClick={() => cleanReq.mutate()}
                    >
                        Clean
                    </Button>
                )}
                <Button
                    sx={SX.button}
                    variant={"contained"}
                    disabled={fetching}
                    onClick={() => setReq.mutate({ ref: refWord, key: keyWord })}
                >
                    Done
                </Button>
            </Box>
        </Grid>
    )
}
