import {useErrorBoundary} from "react-error-boundary"

import {AlertButton} from "../../../shared/component/button/AlertButton"
import {useStoreAction} from "../../../shared/provider/StoreProvider"

export function ClearCache() {
    const {clear} = useStoreAction
    const {resetBoundary} = useErrorBoundary()

    return (
        <AlertButton
            variant={"outlined"}
            label={"Clear"}
            title={"Clear local cache data?"}
            description={`This action will clear all your local cache. It shouldn't cause any difficulties. You will
            lose your active state in Ivory (selection, navigation, some counts will be recalculated). Usually it is
            helpful after updates, when local store was changed and Ivory works incorrectly.`}
            onClick={handleClick}
        />
    )

    function handleClick() {
        clear()
        resetBoundary()
    }
}
