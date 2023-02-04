import {Box, Collapse} from "@mui/material";
import scroll from "../../../style/scroll.module.css";
import {TransitionGroup} from "react-transition-group";
import {CertsItem} from "./CertsItem";
import React from "react";
import {Cert} from "../../../app/types";
import {InfoAlert} from "../../view/InfoAlert";
import {useQuery} from "@tanstack/react-query";
import {certApi} from "../../../app/api";
import {ErrorAlert} from "../../view/ErrorAlert";
import {LinearProgressStateful} from "../../view/LinearProgressStateful";
import {CertTypeProps} from "./CertsContent";

const SX = {
    list: { maxHeight: "500px", overflowY: "auto" },
}

export function CertsList(props: CertTypeProps) {
    const { data: certs, isError, error, isFetching } = useQuery(["certs"], certApi.list)

    return (
        <>
            <LinearProgressStateful color={"inherit"} isFetching={isFetching} line />
            {renderBody()}
        </>
    )

    function renderBody() {
        if (isError) return <ErrorAlert error={error}/>

        const list = Object.entries<Cert>(certs ?? {})
        if (list.length === 0) return <InfoAlert text={"There is no certs yet"}/>

        return (
            <Box sx={SX.list} className={scroll.tiny}>
                <TransitionGroup>
                    {list.map(([key, cert]) => (
                        <Collapse key={key}>
                            <CertsItem uuid={key} cert={cert}/>
                        </Collapse>
                    ))}
                </TransitionGroup>
            </Box>
        )
    }
}
