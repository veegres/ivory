import {useErrorBoundary} from "react-error-boundary";

import {useStoreAction} from "../../../provider/StoreProvider";
import {AlertButton} from "../../view/button/AlertButton";

export function ClearCacheButton() {
    const {clear} = useStoreAction
    const {resetBoundary} = useErrorBoundary()

    return (
        <AlertButton
            variant={"outlined"}
            size={"small"}
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
