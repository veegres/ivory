import {createAsyncThunk, createSlice} from '@reduxjs/toolkit';
import {RootState} from '../../app/store';
import {restClient} from "../../app/hooks";

export interface NodeConfigState {
    data?: object;
    error?: string;
    status: 'idle' | 'loading' | 'failed';
}

const initialState: NodeConfigState = {
    status: 'idle',
};

export const incrementAsync = createAsyncThunk(
    'node/fetchConfig',
    async (node: string, { rejectWithValue }) => {
        try {
            return await restClient<object>(`/node/${node}/config`)
        } catch (error) {
            return rejectWithValue(error.response.data)
        }
    }
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
                state.data = action.payload.response;
            })
            .addCase(incrementAsync.rejected, (state, action) => {
                state.status = 'idle';
                state.error = action.error.message;
            });
    },
});

export const selectNodeConfigData = (state: RootState) => state.nodeConfig.data;
