import {Box} from "@mui/material"

import {SxPropsMap} from "../../helper/HelperType"

const SX: SxPropsMap = {
    greeting: {fontSize: "32px", textAlign: "center"},
}

type Props = {
    username?: string,
}

// PageStartupGreeting is the line that names the person a startup page is
// talking to. It is the most important text on the page - the header above it
// only says which page this is - so it is spelled out once here rather than
// each page picking a heading size of its own.
export function PageStartupGreeting(props: Props) {
    const {username} = props
    return <Box sx={SX.greeting}>Glad to see you, {username}!</Box>
}
