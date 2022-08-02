import {Box, Button, Grid} from "@mui/material";
import React, {ReactNode, useState} from "react";
import {useMutation, useQueryClient} from "react-query";
import {secretApi} from "../../app/api";
import {randomUnicodeAnimal} from "../../app/utils";

const SX = {
    box: { height: "100%", width: "30%", minWidth: "500px" },
    button: { margin: "0 10px" },
    header: { fontSize: '35px', fontWeight: 900, fontFamily: 'monospace', margin: "20px 0", cursor: "pointer" }
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
    const queryClient = useQueryClient();
    const setReq = useMutation(secretApi.set, {
        onSuccess: async () => await queryClient.refetchQueries("secret")
    })
    const cleanReq = useMutation(secretApi.clean, {
        onSuccess: async () => await queryClient.refetchQueries("secret")
    })

    return (
        <Grid container sx={SX.box} direction={"column"} alignItems={"center"} justifyContent={"center"}>
            <Box sx={SX.header} onClick={() => setAnimal(randomUnicodeAnimal())}>
                {header} {animal}
            </Box>
            {children}
            <Box>
                {!clean ? null : (
                    <Button sx={SX.button} variant={"contained"} onClick={() => cleanReq.mutate()}>
                        Clean
                    </Button>
                )}
                <Button sx={SX.button} variant={"contained"} onClick={() => setReq.mutate({ ref: refWord, key: keyWord })}>
                    Done
                </Button>
            </Box>
        </Grid>
    )
}
