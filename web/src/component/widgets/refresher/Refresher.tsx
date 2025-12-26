import {ArrowDropDown, ArrowDropUp, Autorenew} from "@mui/icons-material"
import {Box, Button, ButtonGroup, ClickAwayListener, MenuItem, Paper, Popper, SxProps, Theme} from "@mui/material"
import {QueryKey, useQueryClient} from "@tanstack/react-query"
import {useEffect, useRef, useState} from "react"

import {SxPropsMap} from "../../../app/type"

const SX: SxPropsMap = {
    select: {display: "flex", flexDirection: "column", minWidth: "60px", padding: "5px 0px"},
    item: {padding: "3px 10px", fontSize: "14px", justifyContent: "center"},
    popper: {zIndex: 3},
    icon: {fontSize: "16px", color: "text.secondary"},
    label: {color: "text.secondary", lineHeight: 1},
}

const group = (size: number): SxProps<Theme> => ({"& .MuiButtonGroup-grouped": {
    width: size + 8, height: size, fontSize: "12px", padding: "0px", minWidth: "unset", color: "divider",
}})

const periods: [string, number][] = [
    ["Off", 0], ["1s", 1000], ["5s", 5000], ["10s", 10000],
    ["30s", 30000], ["1m", 60000], ["5m", 50000], ["10m", 100000],
]

type Props = {
    queryKeys: QueryKey[],
    size?: number,
}

export function Refresher(props: Props) {
    const {size = 28, queryKeys} = props
    const queryClient = useQueryClient()
    const [open, setOpen] = useState(false)
    const [period, setPeriod] = useState(periods[0])
    const [label, ms] = period
    const anchorRef = useRef(null)

    // eslint-disable-next-line react-hooks/exhaustive-deps
    useEffect(() => () => {handleClear(periods[0][1])}, [])

    return (
        <>
            <ButtonGroup color={"inherit"} ref={anchorRef} sx={group(size)}>
                <Button onClick={handleRefresh}><Autorenew sx={SX.icon}/></Button>
                <Button onClick={() => setOpen(!open)}>
                    {ms === 0 ? (open ? <ArrowDropUp sx={SX.icon}/> : <ArrowDropDown sx={SX.icon}/>) : <Box sx={SX.label}>{label}</Box>}
                </Button>
            </ButtonGroup>
            { }
            <Popper sx={SX.popper} open={open} anchorEl={anchorRef.current} placement={"bottom"}>
                <Paper>
                    <ClickAwayListener onClickAway={() => setOpen(false)}>
                        <Box sx={SX.select} onClick={() => setOpen(false)}>
                            {periods.map(([label, ms]) => (
                                <MenuItem key={label} sx={SX.item} onClick={() => handleInterval(label, ms)}>
                                    {label}
                                </MenuItem>
                            ))}
                        </Box>
                    </ClickAwayListener>
                </Paper>
            </Popper>
        </>
    )

    function handleClear(ms: number) {
        for (const queryKey of queryKeys) {
            queryClient.setQueryDefaults(queryKey, {refetchInterval: ms})
        }
    }

    async function handleInterval(label: string, ms: number) {
        setPeriod([label, ms])
        for (const queryKey of queryKeys) {
            queryClient.setQueryDefaults(queryKey, {refetchInterval: ms})
            await queryClient.refetchQueries({queryKey})
        }
    }

    async function handleRefresh() {
        for (const queryKey of queryKeys) {
            await queryClient.refetchQueries({queryKey})
        }
    }
}

