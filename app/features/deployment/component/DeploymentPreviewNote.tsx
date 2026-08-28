import {Box} from "@mui/material"

import {CodeToken} from "../../../shared/component/box/CodeToken"
import {Note} from "../../../shared/component/box/Note"
import {SxPropsMap} from "../../../shared/helper/HelperType"
import {DeployPasswordMask} from "../../../shared/helper/HelperUtils"

const SX: SxPropsMap = {
    // NOTE: the inset a CodeField label takes, so the line explaining the
    // fields below starts where their labels do
    box: {padding: "0px 5px"},
}

// DeploymentPreviewNote says what an interpolated command is, on every screen
// that shows one - it covers the post script below it too, which is
// interpolated from the same values. Everything is real except the password, so
// the mask is the only thing left to explain - and it is worth explaining,
// because a reader has to know the deployed command is not the one on screen.
export function DeploymentPreviewNote() {
    return (
        <Box sx={SX.box}>
            <Note>
                {"This is what will run on the node, with "}
                <CodeToken>{DeployPasswordMask}</CodeToken>
                {" standing in for the password: the server substitutes the real one from the database credentials above, so it never reaches the browser."}
            </Note>
        </Box>
    )
}
