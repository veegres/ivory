import {Box, Button, Grid} from "@mui/material";
import React, {ReactNode} from "react";
import {useMutation, useQueryClient} from "react-query";
import {secretApi} from "../../app/api";

const SX = {
    box: { height: "100%", width: "40%", minWidth: "300px" },
    button: { margin: "0 10px" }
}

type Props = {
    children: ReactNode
    keyWord: string
    refWord: string
    clean: boolean
}

export function Secret(props: Props) {
    const { keyWord, refWord, children, clean } = props
    const queryClient = useQueryClient();
    const setReq = useMutation(secretApi.set, {
        onSuccess: async () => await queryClient.refetchQueries("secret")
    })
    const cleanReq = useMutation(secretApi.clean, {
        onSuccess: async () => await queryClient.refetchQueries("secret")
    })

    return (
        <Grid container sx={SX.box} direction={"column"} alignItems={"center"} justifyContent={"center"}>
            {children}
            <Box>
                <Button sx={SX.button} variant={"contained"} onClick={() => setReq.mutate({ ref: refWord, key: keyWord })}>
                    Done
                </Button>
                {!clean ? null : (
                    <Button sx={SX.button} variant={"contained"} onClick={() => cleanReq.mutate()}>
                        Clean
                    </Button>
                )}
            </Box>
        </Grid>
    )
}
