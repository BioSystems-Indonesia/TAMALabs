import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Chip,
  Collapse,
  IconButton,
  Paper,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  InputAdornment,
  Autocomplete
} from '@mui/material';
import {
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  Add as AddIcon,
  Settings as SettingsIcon,
  Search as SearchIcon,
} from '@mui/icons-material';
import FileDownloadIcon from '@mui/icons-material/FileDownload';
import useAxios from '../../hooks/useAxios';
import { useNotify } from 'react-admin';

// Define the structure of a log entry from the backend
interface LogEntry {
  time: string;
  level: string;
  msg: string;
  uri: string;
  error: string;
  [key: string]: string | number | boolean | object | null; // Allow for additional JSON fields
}

interface CustomColumn {
  field: string;
  headerName: string;
  width?: number;
}

const LogViewer = () => {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());
  const [customColumns, setCustomColumns] = useState<CustomColumn[]>([]);
  const [showColumnDialog, setShowColumnDialog] = useState(false);
  const [newColumn, setNewColumn] = useState<Partial<CustomColumn>>({});
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(50);
  const [searchTerm, setSearchTerm] = useState('');
  const axios = useAxios()
  const notify = useNotify()

  const exportLogs = () => {
    axios({
      url: "/log/export",
      method: "GET",
      responseType: "blob",
    }).then((res) => {
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", "logs.zip");
      document.body.appendChild(link);
      link.click();
    }).catch((err) => {
      notify("Error export logs: " + err, {
        type: 'error',
      });
    })
  }

  useEffect(() => {
    // Set up global function for JSON field clicks
    (window as { addColumnFromJson?: (fieldName: string) => void }).addColumnFromJson = addColumnFromJson;

    // 1. Connect to the SSE endpoint from our Go backend
    const baseUrl = import.meta.env.VITE_BACKEND_BASE_URL;
    const eventSource = new EventSource(`${baseUrl}/log/stream`);

    // 2. Handle incoming messages
    eventSource.onmessage = (event) => {
      const newLog = JSON.parse(event.data);
      console.log(newLog);
      // Add the new log to the top of the list
      setLogs((prevLogs) => [newLog, ...prevLogs]);
    };

    // Handle any errors
    eventSource.onerror = (err) => {
      console.error('EventSource failed:', err);
      eventSource.close();
    };

    // 3. Clean up the connection when the component unmounts
    return () => {
      console.log('unmount');
      eventSource.close();
      // Clean up global function
      delete (window as { addColumnFromJson?: (fieldName: string) => void }).addColumnFromJson;
    };
  }, []); // Empty dependency array means this effect runs once on mount

  const handleRowClick = (rowId: string) => {
    setExpandedRows(prev => {
      const newSet = new Set(prev);
      if (newSet.has(rowId)) {
        newSet.delete(rowId);
      } else {
        newSet.add(rowId);
      }
      return newSet;
    });
  };

  const getLevelColor = (level: string) => {
    switch (level.toLowerCase()) {
      case 'error':
        return 'error';
      case 'warn':
      case 'warning':
        return 'warning';
      case 'info':
        return 'info';
      case 'debug':
        return 'default';
      default:
        return 'default';
    }
  };


  const formatJsonWithClickableFields = (obj: unknown, onFieldClick: (fieldName: string) => void) => {
    try {
      const jsonStr = JSON.stringify(obj, null, 2);
      return jsonStr.replace(/"([^"]+)":/g, (match, fieldName) => {
        if (fieldName !== 'time' && fieldName !== 'level' && fieldName !== 'msg') {
          return `<span class="json-field" onclick="window.addColumnFromJson('${fieldName}')">"${fieldName}"</span>:`;
        }
        return match;
      });
    } catch {
      return String(obj);
    }
  };

  const addCustomColumn = () => {
    if (newColumn.field) {
      const headerName = newColumn.field
        .split('_')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ');

      setCustomColumns(prev => [...prev, {
        field: newColumn.field!,
        headerName: headerName,
        width: newColumn.width
      }]);
      setNewColumn({});
      setShowColumnDialog(false);
    }
  };

  const removeCustomColumn = (field: string) => {
    setCustomColumns(prev => prev.filter(col => col.field !== field));
  };

  // Get all available field names from logs
  const getAvailableFields = (): string[] => {
    const fields = new Set<string>();
    logs.forEach(log => {
      Object.keys(log).forEach(key => {
        if (key !== 'time' && key !== 'level' && key !== 'msg') {
          fields.add(key);
        }
      });
    });
    return Array.from(fields).sort();
  };

  // Add column from JSON click
  const addColumnFromJson = (fieldName: string) => {
    const headerName = fieldName
      .split('_')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ');

    setCustomColumns(prev => [...prev, {
      field: fieldName,
      headerName: headerName
    }]);
  };

  // Filter logs based on search term
  const filteredLogs = logs.filter(log =>
    log.msg.toLowerCase().includes(searchTerm.toLowerCase()) ||
    log.level.toLowerCase().includes(searchTerm.toLowerCase()) ||
    log.time.toLowerCase().includes(searchTerm.toLowerCase()) ||
    customColumns.some(col => {
      const value = log[col.field];
      return value && String(value).toLowerCase().includes(searchTerm.toLowerCase());
    })
  );

  // Paginate logs
  const paginatedLogs = filteredLogs.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);

  const handleChangePage = (_event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  return (
    <Box sx={{ height: '80vh', width: '100%' }}>
      {/* Header with controls */}
      <Box sx={{ mb: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h6">Log Viewer</Typography>
        <Stack direction="row" spacing={1}>
          <TextField
            size="small"
            placeholder="Search logs..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
            sx={{ width: 250 }}
          />
          <Button
            variant="outlined"
            startIcon={<FileDownloadIcon />}
            onClick={() => exportLogs()}
          >
            Export Logs
          </Button>
          <Button
            variant="outlined"
            startIcon={<AddIcon />}
            onClick={() => setShowColumnDialog(true)}
          >
            Add Column
          </Button>
          {customColumns.length > 0 && (
            <Button
              variant="outlined"
              startIcon={<SettingsIcon />}
              onClick={() => setShowColumnDialog(true)}
            >
              Manage Columns ({customColumns.length})
            </Button>
          )}
        </Stack>
      </Box>

      {/* Custom columns display */}
      {customColumns.length > 0 && (
        <Box sx={{ mb: 2 }}>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
            Custom Columns:
          </Typography>
          <Stack direction="row" spacing={1} flexWrap="wrap">
            {customColumns.map((col) => (
              <Chip
                key={col.field}
                label={`${col.headerName} (${col.field})`}
                onDelete={() => removeCustomColumn(col.field)}
                size="small"
                variant="outlined"
              />
            ))}
          </Stack>
        </Box>
      )}

      {/* Table */}
      <Paper sx={{ height: 'calc(100%)', overflow: 'auto', maxWidth: 'calc(95vw - 200px)' }}>
        <TableContainer>
          <Table stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell sx={{ minWidth: 150 }}>Time</TableCell>
                <TableCell sx={{ minWidth: 100 }}>Level</TableCell>
                <TableCell sx={{ minWidth: 300 }}>Message</TableCell>
                <TableCell sx={{ minWidth: 300 }}>Error</TableCell>
                {customColumns.map((col) => (
                  <TableCell key={col.field} sx={{ minWidth: col.width || 150 }}>
                    {col.headerName}
                  </TableCell>
                ))}
                <TableCell sx={{ minWidth: 80 }}>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {paginatedLogs.map((log, index) => {
                const rowId = `${page * rowsPerPage + index}`;
                const isExpanded = expandedRows.has(rowId);

                return (
                  <React.Fragment key={rowId}>
                    <TableRow
                      hover
                      onClick={() => handleRowClick(rowId)}
                      sx={{ cursor: 'pointer' }}
                    >
                      <TableCell>
                        {new Date(log.time).toLocaleTimeString()}
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={log.level}
                          color={getLevelColor(log.level) as 'error' | 'warning' | 'info' | 'default'}
                          size="small"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>
                        <Typography
                          variant="body2"
                          sx={{
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                            maxWidth: '100%'
                          }}
                        >
                          {log.msg}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Typography
                          variant="body2"
                          sx={{
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                            maxWidth: '100%'
                          }}
                        >
                          {log.error}
                        </Typography>
                      </TableCell>


                      {customColumns.map((col) => {
                        const value = log[col.field];
                        return (
                          <TableCell key={col.field}>
                            {value ? String(value) : '-'}
                          </TableCell>
                        );
                      })}
                      <TableCell>
                        <IconButton size="small">
                          {isExpanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                        </IconButton>
                      </TableCell>
                    </TableRow>

                    {/* Expanded row content */}
                    <TableRow>
                      <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={3 + customColumns.length + 1}>
                        <Collapse in={isExpanded} timeout="auto" unmountOnExit>
                          <Box sx={{ margin: 1 }}>
                            <Typography variant="h6" gutterBottom>
                              Full Log Entry
                            </Typography>
                            <Box
                              component="pre"
                              sx={{
                                backgroundColor: 'grey.100',
                                p: 2,
                                borderRadius: 1,
                                overflow: 'auto',
                                fontSize: '0.875rem',
                                fontFamily: 'monospace',
                                maxHeight: 300,
                                '& .json-field': {
                                  cursor: 'pointer',
                                  color: '#1976d2',
                                  textDecoration: 'underline',
                                  '&:hover': {
                                    backgroundColor: 'rgba(25, 118, 210, 0.1)',
                                  }
                                }
                              }}
                              dangerouslySetInnerHTML={{
                                __html: formatJsonWithClickableFields(log, addColumnFromJson)
                              }}
                            />
                          </Box>
                        </Collapse>
                      </TableCell>
                    </TableRow>
                  </React.Fragment>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>

        <TablePagination
          rowsPerPageOptions={[25, 50, 100]}
          component="div"
          count={filteredLogs.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </Paper>

      {/* Add Column Dialog */}
      <Dialog open={showColumnDialog} onClose={() => setShowColumnDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Add Custom Column</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 1 }}>
            <Autocomplete
              freeSolo
              options={getAvailableFields()}
              value={newColumn.field || ''}
              onChange={(_event, newValue) => setNewColumn(prev => ({ ...prev, field: newValue || '' }))}
              onInputChange={(_event, newInputValue) => setNewColumn(prev => ({ ...prev, field: newInputValue }))}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Field Name"
                  placeholder="e.g., user_id, request_id"
                  helperText="Type or select a field name from existing logs"
                />
              )}
            />
            <TextField
              label="Width (optional)"
              type="number"
              value={newColumn.width || ''}
              onChange={(e) => setNewColumn(prev => ({ ...prev, width: parseInt(e.target.value) || undefined }))}
              placeholder="150"
              helperText="Column width in pixels"
            />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowColumnDialog(false)}>Cancel</Button>
          <Button onClick={addCustomColumn} variant="contained" disabled={!newColumn.field}>
            Add Column
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default LogViewer;