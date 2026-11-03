import { Cancel, CheckCircle, Delete, Edit } from "@mui/icons-material";
import { Box, Button, Chip, TextField, Typography } from "@mui/material";
import { useState } from "react";

import { ErrorLeaderMissing } from "../../../shared/component/box/ErrorManual";
import { ErrorSmart } from "../../../shared/component/box/ErrorSmart";
import { NoBox } from "../../../shared/component/box/NoBox";
import { TitleBox } from "../../../shared/component/box/TitleBox";
import { SimpleButton } from "../../../shared/component/button/SimpleButton";
import { FieldRow } from "../../../shared/component/input/FieldRow";
import { SkeletonGroup } from "../../../shared/component/progress/SkeletonGroup";
import { Feature } from "../../Feature";
import { ManageAccess } from "../../management/component/ManageAccess";
import { Role } from "../../node/api/NodeType";
import {
  useRouterDbDatabaseAdminCreate,
  useRouterDbDatabaseAdminDelete,
  useRouterDbDatabaseAdminList,
  useRouterDbDatabaseAdminUpdate,
} from "../api/QueryHook";
import { Connection, DbDatabaseItem } from "../api/QueryType";
import { QueryDbAdminPanel, QueryDbAdminSX } from "./QueryDbAdminPanel";

type Props = {
  connection: Connection;
  role: Role;
};

export function QueryDbDatabaseList(props: Props) {
  const { connection, role } = props;
  const list = useRouterDbDatabaseAdminList(connection);
  const [newName, setNewName] = useState("");
  const create = useRouterDbDatabaseAdminCreate(connection, () =>
    setNewName(""),
  );

  if (role !== "leader") return <ErrorLeaderMissing />;

  return (
    <QueryDbAdminPanel
      connection={connection}
      viewFeature={Feature.ViewQueryDbDatabaseList}
      createFeature={Feature.ManageQueryDbDatabaseCreate}
      listKind={"database"}
      createLabel={"New database"}
      alert={
        "User databases on this instance. System catalogs (postgres, template0, template1) are listed but cannot be changed here."
      }
      renderCreate={renderCreateForm()}
    >
      {renderList()}
    </QueryDbAdminPanel>
  );

  function renderCreateForm() {
    return (
      <>
        <FieldRow>
          <TextField
            size={"small"}
            label={"Name"}
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
          />
        </FieldRow>
        <Button
          sx={{ marginTop: 1 }}
          size={"small"}
          variant={"contained"}
          color={"secondary"}
          disabled={!newName.trim() || create.isPending}
          onClick={() => create.mutate({ connection, name: newName.trim() })}
        >
          Create database
        </Button>
      </>
    );
  }

  function renderList() {
    if (list.isPending)
      return (
        <SkeletonGroup count={4} width={"100%"} height={72} grow={false} />
      );
    if (list.isError) return <ErrorSmart error={list.error} />;
    if (!list.data || list.data.length === 0)
      return <NoBox text={"There are no databases yet"} />;

    return list.data.map((item) => (
      <QueryDbDatabaseRow
        key={item.name}
        connection={connection}
        database={item}
      />
    ));
  }
}

type RowProps = {
  connection: Connection;
  database: DbDatabaseItem;
};

function QueryDbDatabaseRow(props: RowProps) {
  const { connection, database } = props;
  const [edit, setEdit] = useState(false);
  const [newName, setNewName] = useState("");
  const protectedDb = isProtectedDatabase(database.name);

  const rename = useRouterDbDatabaseAdminUpdate(connection, () => {
    setEdit(false);
    setNewName("");
  });
  const remove = useRouterDbDatabaseAdminDelete(connection);

  return (
    <TitleBox
      label={database.name}
      island={true}
      dense={true}
      collapsible={false}
    >
      <Box sx={QueryDbAdminSX.cardHead}>
        <Box sx={QueryDbAdminSX.metaGrid}>
          <Box>
            <Typography sx={QueryDbAdminSX.metaLabel}>Owner</Typography>
            <Typography sx={QueryDbAdminSX.metaValue}>
              {database.owner}
            </Typography>
          </Box>
          <Box>
            <Typography sx={QueryDbAdminSX.metaLabel}>Size</Typography>
            <Typography sx={QueryDbAdminSX.metaValue}>
              {database.size}
            </Typography>
          </Box>
          <Box>
            <Typography sx={QueryDbAdminSX.metaLabel}>Encoding</Typography>
            <Typography sx={QueryDbAdminSX.metaValue}>
              {database.encoding}
            </Typography>
          </Box>
        </Box>
        {protectedDb ? (
          <Chip size={"small"} label={"protected"} variant={"outlined"} />
        ) : !edit ? (
          renderActions()
        ) : undefined}
      </Box>
      {edit && (
        <FieldRow>
          <TextField
            size={"small"}
            label={"New name"}
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
          />
          <SimpleButton
            tooltip={"Cancel"}
            onClick={() => setEdit(false)}
            disabled={rename.isPending}
          >
            <Cancel fontSize={"small"} />
          </SimpleButton>
          <SimpleButton
            tooltip={"Rename"}
            loading={rename.isPending}
            disabled={!newName.trim() || newName.trim() === database.name}
            onClick={() =>
              rename.mutate({
                connection,
                name: database.name,
                newName: newName.trim(),
              })
            }
          >
            <CheckCircle fontSize={"small"} />
          </SimpleButton>
        </FieldRow>
      )}
    </TitleBox>
  );

  function renderActions() {
    return (
      <Box sx={QueryDbAdminSX.cardActions}>
        <ManageAccess feature={Feature.ManageQueryDbDatabaseUpdate}>
          <SimpleButton
            tooltip={"Rename"}
            onClick={() => setEdit(true)}
            disabled={remove.isPending}
          >
            <Edit fontSize={"small"} />
          </SimpleButton>
        </ManageAccess>
        <ManageAccess feature={Feature.ManageQueryDbDatabaseDelete}>
          <SimpleButton
            tooltip={"Delete database"}
            loading={remove.isPending}
            onClick={() => remove.mutate({ connection, name: database.name })}
          >
            <Delete fontSize={"small"} />
          </SimpleButton>
        </ManageAccess>
      </Box>
    );
  }
}

function isProtectedDatabase(name: string) {
  return name === "postgres" || name === "template0" || name === "template1";
}
