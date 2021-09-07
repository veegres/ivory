import {createAsyncThunk, createSlice} from '@reduxjs/toolkit';
import {RootState} from '../../app/store';
import {restClient} from "../../app/hooks";

export interface NodeConfigState {
    data?: object;
    status: 'idle' | 'loading' | 'failed';
}

const initialState: NodeConfigState = {
    data: undefined,
    status: 'idle',
};

export const incrementAsync = createAsyncThunk(
    'node/fetchConfig',
    async (node: string) => await restClient<object>(`/node/${node}/config`)
);

export const nodeConfigSlice = createSlice({
    name: 'node/config',
    initialState,
    reducers: {},
    extraReducers: (builder) => {
        builder
            .addCase(incrementAsync.pending, (state) => {
                state.status = 'loading';
            })
            .addCase(incrementAsync.fulfilled, (state, action) => {
                state.status = 'idle';
                state.data = action.payload;
            });
    },
});

export const selectNodeConfigData = (state: RootState) => state.nodeConfig.data;
