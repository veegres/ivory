import {Collapse} from "@mui/material";
import {TransitionGroup} from "react-transition-group";
import {CertsItem} from "./CertsItem";
import React from "react";
import {Cert} from "../../../app/types";
import {InfoAlert} from "../../view/InfoAlert";
import {useQuery} from "@tanstack/react-query";
import {certApi} from "../../../app/api";
import {ErrorAlert} from "../../view/ErrorAlert";
import {LinearProgressStateful} from "../../view/LinearProgressStateful";
import {CertTypeProps} from "./Certs";
import {MenuWrapperScroll} from "../menu/MenuWrapperScroll";

export function CertsList(props: CertTypeProps) {
    const {type} = props
    const query = useQuery(["certs", type], () => certApi.list(type))
    const {data: certs, isError, error, isFetching} = query

    return (
        <>
            <LinearProgressStateful color={"inherit"} isFetching={isFetching} line/>
            {renderBody()}
        </>
    )

    function renderBody() {
        if (isError) return <ErrorAlert error={error}/>

        const list = Object.entries<Cert>(certs ?? {})
        if (list.length === 0) return <InfoAlert text={"There is no files yet"}/>

        return (
            <MenuWrapperScroll>
                <TransitionGroup>
                    {list.map(([key, cert]) => (
                        <Collapse key={key}>
                            <CertsItem uuid={key} cert={cert}/>
                        </Collapse>
                    ))}
                </TransitionGroup>
            </MenuWrapperScroll>
        )
    }
}
