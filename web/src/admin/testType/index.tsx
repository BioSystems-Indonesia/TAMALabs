import {
  Stack,
  useTheme,
  Card,
  CardContent,
  Chip,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from "@mui/material";
import Box from "@mui/material/Box";
import CreateNewFolderIcon from "@mui/icons-material/CreateNewFolder";
import { useTheme as useMuiTheme } from "@mui/material/styles";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState, useRef } from "react";
import {
  ArrayInput,
  AutocompleteInput,
  BooleanField,
  BooleanInput,
  Button,
  Create,
  CreateButton,
  Datagrid,
  Edit,
  FunctionField,
  List,
  NumberInput,
  Show,
  SimpleForm,
  SimpleFormIterator,
  TextField,
  TextInput,
  TopToolbar,
  required,
} from "react-admin";
import { useFormContext } from "react-hook-form";
import { useSearchParams } from "react-router-dom";
import FeatureList from "../../component/FeatureList";
import type { ActionKeys } from "../../types/props";
import type { Unit } from "../../types/unit";
import { TestFilterSidebar } from "../workOrder/TestTypeFilter";
import useAxios from "../../hooks/useAxios";
import { TestType } from "../../types/test_type";
import { Device } from "../../types/device";
import type { ObservationRequestCreateRequest } from "../../types/observation_requests";

const NullableField = ({ value }: { value: any }) => {
  // Check if value is null or undefined, but allow 0
  const isNull = value === null || value === undefined || value === '';

  return (
    <span style={{
      color: isNull ? '#888' : 'inherit',
      fontStyle: isNull ? 'italic' : 'normal',
      opacity: isNull ? 0.6 : 1,
      fontSize: isNull ? '0.875rem' : 'inherit'
    }}>
      {isNull ? 'null' : value}
    </span>
  );
};

const SpecificReferenceRangesDisplay = () => {
  const { watch } = useFormContext();
  const theme = useTheme();
  const specificRanges = watch('specific_ref_ranges');

  if (!specificRanges || specificRanges.length === 0) {
    return (
      <Box sx={{ mt: 2 }}>
        <Typography
          variant="body2"
          sx={{
            color: theme.palette.text.secondary,
            fontStyle: "italic",
          }}
        >
          No specific reference ranges defined. Using default ranges for all patients.
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ mt: 3 }}>
      <Typography
        variant="body1"
        sx={{
          fontWeight: 500,
          color: theme.palette.text.secondary,
          mb: 2,
        }}
      >
        🎯 Specific Reference Ranges
      </Typography>
      <TableContainer component={Paper} variant="outlined">
        <Table size="small">
          <TableHead>
            <TableRow sx={{ backgroundColor: theme.palette.action.hover }}>
              <TableCell><strong>Gender</strong></TableCell>
              <TableCell><strong>Age Range</strong></TableCell>
              <TableCell><strong>Reference Range</strong></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {specificRanges.map((range: any, index: number) => {
              const genderLabel = !range.gender
                ? "All"
                : range.gender === "M"
                  ? "Male"
                  : "Female";

              let ageLabel = "All Ages";
              if (range.age_min && range.age_max) {
                ageLabel = `${range.age_min} - ${range.age_max} years`;
              } else if (range.age_min) {
                ageLabel = `≥ ${range.age_min} years`;
              } else if (range.age_max) {
                ageLabel = `≤ ${range.age_max} years`;
              }

              let rangeLabel = "-";
              if (range.normal_ref_string) {
                rangeLabel = range.normal_ref_string;
              } else if (
                range.low_ref_range !== null &&
                range.low_ref_range !== undefined &&
                range.high_ref_range !== null &&
                range.high_ref_range !== undefined
              ) {
                rangeLabel = `${range.low_ref_range} - ${range.high_ref_range}`;
              }

              return (
                <TableRow key={index}>
                  <TableCell>
                    <Chip
                      label={genderLabel}
                      size="small"
                      color={
                        range.gender === "M"
                          ? "primary"
                          : range.gender === "F"
                            ? "secondary"
                            : "default"
                      }
                    />
                  </TableCell>
                  <TableCell>{ageLabel}</TableCell>
                  <TableCell>
                    <strong>{rangeLabel}</strong>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
};

export const TestTypeDatagrid = () => {
  return (
    <Datagrid bulkActionButtons={false}>
      <TextField source="id" />
      <FunctionField
        label="Name"
        render={(record: TestType) => <NullableField value={record.name} />}
      />
      <TextField source="code" />
      <FunctionField
        label="Alternative Codes"
        render={(record: TestType) => {
          if (record.alternative_codes && record.alternative_codes.length > 0) {
            return <NullableField value={record.alternative_codes.join(', ')} />;
          }
          return <NullableField value={null} />;
        }}
      />
      <FunctionField
        label="Category"
        render={(record: TestType) => <NullableField value={record.category} />}
      />
      <FunctionField
        label="Sub Category"
        render={(record: TestType) => <NullableField value={record.sub_category} />}
      />
      <FunctionField
        label="Low"
        render={(record: TestType) => <NullableField value={record.low_ref_range} />}
      />
      <FunctionField
        label="High"
        render={(record: TestType) => <NullableField value={record.high_ref_range} />}
      />
      <FunctionField
        label="Normal String"
        render={(record: TestType) => <NullableField value={record.normal_ref_string} />}
      />
      <BooleanField source="is_calculated_test" label="Calc Test" sortable />
      <FunctionField
        label="Unit"
        render={(record: TestType) => <NullableField value={record.unit} />}
      />
      <FunctionField
        label="Devices"
        render={(record: TestType) => {
          if (record.devices && record.devices.length > 0) {
            return <NullableField value={record.devices.map(d => d.name).join(', ')} />;
          }
          // Fallback to old device field for backward compatibility
          if (record.device) {
            return <NullableField value={record.device.name} />;
          }
          return <NullableField value="General" />;
        }}
      />
      <FunctionField
        label="Types"
        render={(record: TestType) =>
          <NullableField
            value={
              record.types && record.types.length > 0
                ? record.types.map((t) => t.type).join(", ")
                : null
            }
          />
        }
      />
      <FunctionField
        label="Decimal"
        render={(record: TestType) => <NullableField value={record.decimal} />}
      />
    </Datagrid>
  );
};

export function TestTypeShow() {
  const theme = useMuiTheme();

  return (
    <Box
      sx={{
        minHeight: "100vh",
        bgcolor: theme.palette.background.default,
        pb: 4,
      }}
    >
      <Show resource="test-type">
        <TestTypeForm readonly mode={"SHOW"} />
      </Show>
    </Box>
  );
}

export const TestTypeList = () => {
  const [selectedData, setSelectedData] = useState<Record<number, ObservationRequestCreateRequest>>({});

  return (
    <List
      actions={<TestTypeListActions />}
      aside={
        <TestFilterSidebar
          selectedData={selectedData}
          setSelectedData={setSelectedData}
        />
      }
      title="Test Type"
      sort={{
        field: "id",
        order: "DESC",
      }}
      perPage={50}
      sx={{
        "& .RaList-main": {},
        "& .RaList-content": {
          backgroundColor: "background.paper",
          padding: 2,
          borderRadius: 1,
        },
      }}
      storeKey={false}
      exporter={false}
    >
      <TestTypeDatagrid />
    </List>
  );
};

function ReferenceSection() {
  // This component can be used to show additional reference range information
  // Currently integrated directly in the form
  return null;
}

type TestTypeFormProps = {
  readonly?: boolean;
  mode?: ActionKeys;
};

function TestTypeInput(props: TestTypeFormProps) {
  const theme = useTheme();
  const axios = useAxios();
  const { data: filter, isLoading: isFilterLoading } = useQuery({
    queryKey: ["filterTestType"],
    queryFn: () => axios.get("/test-type/filter").then((res) => res.data),
  });

  const [categories, setCategories] = useState<string[]>([]);
  const [subCategories, setSubCategories] = useState<string[]>([]);

  useEffect(() => {
    if (filter) {
      setCategories(filter.categories);
      setSubCategories(filter.sub_categories);
    }
  }, [filter, isFilterLoading]);

  const { data: units, isLoading: isUnitLoading } = useQuery<Unit[]>({
    queryKey: ["unit"],
    queryFn: () => axios.get("/unit").then((res) => res.data),
  });

  const { data: devices, isLoading: isDeviceLoading } = useQuery({
    queryKey: ["devices"],
    queryFn: () => axios.get("/device").then((res) => res.data),
  });

  const [unit, setUnit] = useState<string[]>([]);
  useEffect(() => {
    if (units && Array.isArray(units)) {
      const unitValues = units.map((unit) => unit.value);
      const uniqueUnits = [...new Set(unitValues)];
      setUnit(uniqueUnits);
    }
  }, [units, isUnitLoading]);

  const { setValue, watch } = useFormContext();
  const [params] = useSearchParams();

  // Watch for devices and set device_ids when component mounts or devices change
  const devicesValue = watch('devices');
  const deviceIdsValue = watch('device_ids');

  useEffect(() => {
    // If devices exist but device_ids doesn't, populate device_ids
    if (devicesValue && Array.isArray(devicesValue) && devicesValue.length > 0) {
      if (!deviceIdsValue || deviceIdsValue.length === 0) {
        const ids = devicesValue.map((d: any) => d.id);
        console.log('Setting device_ids from devices:', ids);
        setValue('device_ids', ids);
      }
    }
  }, [devicesValue, deviceIdsValue, setValue]);
  useEffect(() => {
    const hasCodeParam = params.has("code");
    if (hasCodeParam) {
      const code = params.get("code");
      setValue("code", code);
      setValue("name", code);
    }
  }, [params, setValue]);

  return (
    <Stack spacing={3} sx={{ width: "100%" }}>
      <Card
        elevation={0}
        sx={{
          border: `1px solid ${theme.palette.divider}`,
          borderRadius: 2,
        }}
      >
        <CardContent sx={{ p: 3 }}>
          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              gap: 1.5,
              mb: 3,
            }}
          >
            <Typography
              variant="subtitle1"
              sx={{
                fontWeight: 600,
                color: theme.palette.text.primary,
              }}
            >
              ❗Basic Information
            </Typography>
            <Chip
              label="Required"
              size="small"
              color="error"
              variant="outlined"
              sx={{ ml: "auto", fontSize: "0.75rem" }}
            />
          </Box>

          <Stack spacing={3}>
            <Stack direction={"row"} gap={3} width={"100%"}>
              <TextInput
                source="name"
                readOnly={props.readonly}
                validate={[required()]}
                fullWidth
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 2,
                    transition: "all 0.2s ease",
                    ...(!props.readonly && {
                      "&:hover": {
                        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                      },
                    }),
                  },
                }}
              />
              <TextInput
                source="code"
                readOnly={props.readonly}
                validate={[required()]}
                fullWidth
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 2,
                    transition: "all 0.2s ease",
                    ...(!props.readonly && {
                      "&:hover": {
                        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                      },
                    }),
                  },
                }}
              />

              <TextInput
                source="alias_code"
                label="Alias Code (SIMRS)"
                helperText="Optional: For SIMRS integration"
                readOnly={props.readonly}
                fullWidth
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 2,
                    transition: "all 0.2s ease",
                    ...(!props.readonly && {
                      "&:hover": {
                        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                      },
                    }),
                  },
                }}
              />
            </Stack>

            <Stack direction={"row"} gap={3} width={"100%"}>
              <TextInput
                source="loinc_code"
                label="LOINC Code"
                helperText="Optional: Logical Observation Identifiers Names and Codes"
                readOnly={props.readonly}
                fullWidth
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 2,
                    transition: "all 0.2s ease",
                    ...(!props.readonly && {
                      "&:hover": {
                        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                      },
                    }),
                  },
                }}
              />
            </Stack>

            <Box>
              <Typography
                variant="body2"
                sx={{
                  fontWeight: 500,
                  color: theme.palette.text.secondary,
                  mb: 2,
                  display: "flex",
                  alignItems: "center",
                  gap: 1,
                }}
              >
                🔄 Alternative Codes
                <Chip
                  label="Optional"
                  size="small"
                  color="info"
                  variant="outlined"
                  sx={{ fontSize: "0.7rem" }}
                />
              </Typography>
              <Typography
                variant="caption"
                sx={{
                  color: theme.palette.text.secondary,
                  display: "block",
                  mb: 2,
                }}
              >
                Add alternative codes from different devices that map to the same test. For example, if Device A sends "HB" and Device B sends "HEMO", both will be recognized as the same test.
              </Typography>

              <ArrayInput source="alternative_codes">
                <SimpleFormIterator
                  inline
                  disableReordering={props.readonly}
                  disableAdd={props.readonly}
                  disableRemove={props.readonly}
                  sx={{
                    "& .RaSimpleFormIterator-line": {
                      mb: 1,
                    },
                  }}
                >
                  <TextInput
                    source=""
                    label="Code"
                    helperText="e.g., HB, HEMO, Hgb"
                    readOnly={props.readonly}
                    sx={{
                      "& .MuiOutlinedInput-root": {
                        borderRadius: 1.5,
                        minWidth: "200px",
                      },
                    }}
                  />
                </SimpleFormIterator>
              </ArrayInput>
            </Box>

            <Stack direction={"row"} gap={3} width={"100%"}>
              <AutocompleteInput
                source="category"
                readOnly={props.readonly}
                filterSelectedOptions={false}
                loading={isFilterLoading}
                choices={categories.map((val) => ({ id: val, name: val }))}
                onCreate={(val) => {
                  if (!val || categories.includes(val)) return;
                  const newCategories = [...categories, val];
                  setCategories(newCategories);
                  return { id: val, name: val };
                }}
                fullWidth
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 2,
                    transition: "all 0.2s ease",
                    ...(!props.readonly && {
                      "&:hover": {
                        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                      },
                    }),
                  },
                }}
              />
              <AutocompleteInput
                source="sub_category"
                readOnly={props.readonly}
                loading={isFilterLoading}
                choices={subCategories.map((val) => ({ id: val, name: val }))}
                onCreate={(val) => {
                  if (!val || subCategories.includes(val)) return;
                  const newSubCategories = [...subCategories, val];
                  setSubCategories(newSubCategories);
                  return { id: val, name: val };
                }}
                fullWidth
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 2,
                    transition: "all 0.2s ease",
                    ...(!props.readonly && {
                      "&:hover": {
                        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                      },
                    }),
                  },
                }}
              />
            </Stack>
          </Stack>
        </CardContent>
      </Card>

      {/* Device Selection Card */}
      <Card
        elevation={0}
        sx={{
          border: `1px solid ${theme.palette.divider}`,
          borderRadius: 2,
        }}
      >
        <CardContent sx={{ p: 3 }}>
          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              gap: 1.5,
              mb: 3,
            }}
          >
            <Typography
              variant="subtitle1"
              sx={{
                fontWeight: 600,
                color: theme.palette.text.primary,
              }}
            >
              🔬 Device Assignment
            </Typography>
            <Chip
              label="Optional"
              size="small"
              color="primary"
              variant="outlined"
              sx={{ ml: "auto", fontSize: "0.75rem" }}
            />
          </Box>

          <Stack spacing={3}>
            <AutocompleteInput
              source="device_ids"
              label="Assigned Devices"
              helperText="Select one or more devices that can perform this test. Multiple selection is supported."
              readOnly={props.readonly}
              loading={isDeviceLoading}
              multiple
              choices={
                devices && Array.isArray(devices)
                  ? devices.map((device: Device) => ({
                    id: device.id,
                    name: `${device.name} (${device.type})`,
                  }))
                  : []
              }
              format={(value: any) => {
                console.log('AutocompleteInput format value:', value);
                return value || [];
              }}
              parse={(value: any) => {
                console.log('AutocompleteInput parse value:', value);
                return value;
              }}
              fullWidth
              sx={{
                "& .MuiOutlinedInput-root": {
                  borderRadius: 2,
                  transition: "all 0.2s ease",
                  ...(!props.readonly && {
                    "&:hover": {
                      boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                    },
                  }),
                },
              }}
            />
          </Stack>
        </CardContent>
      </Card>

      <Card
        elevation={0}
        sx={{
          border: `1px solid ${theme.palette.divider}`,
          borderRadius: 2,
        }}
      >
        <CardContent sx={{ p: 3 }}>
          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              gap: 1.5,
              mb: 3,
            }}
          >
            <Typography
              variant="subtitle1"
              sx={{
                fontWeight: 600,
                color: theme.palette.text.primary,
              }}
            >
              📊 Range & Units
            </Typography>
            <Chip
              label="Required"
              size="small"
              color="error"
              variant="outlined"
              sx={{ ml: "auto", fontSize: "0.75rem" }}
            />
          </Box>

          <Stack spacing={3}>
            <Stack direction={"row"} gap={3} width={"100%"}>
              <NumberInput
                source="low_ref_range"
                label="Low Range"
                readOnly={props.readonly}
                validate={[required()]}
                fullWidth
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 2,
                    transition: "all 0.2s ease",
                    ...(!props.readonly && {
                      "&:hover": {
                        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                      },
                    }),
                  },
                }}
              />
              <NumberInput
                source="high_ref_range"
                label="High Range"
                readOnly={props.readonly}
                validate={[required()]}
                fullWidth
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 2,
                    transition: "all 0.2s ease",
                    ...(!props.readonly && {
                      "&:hover": {
                        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                      },
                    }),
                  },
                }}
              />
              <AutocompleteInput
                source="unit"
                readOnly={props.readonly}
                loading={isUnitLoading}
                choices={[...new Set(unit)].map((val) => ({
                  id: val,
                  name: val,
                }))}
                onCreate={(val) => {
                  if (!val || unit.includes(val)) return;
                  const newUnit = [...new Set([...unit, val])];
                  setUnit(newUnit);
                  return { id: val, name: val };
                }}
                fullWidth
                sx={{
                  "& .MuiOutlinedInput-root": {
                    borderRadius: 2,
                    transition: "all 0.2s ease",
                    ...(!props.readonly && {
                      "&:hover": {
                        boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                      },
                    }),
                  },
                }}
              />
            </Stack>
            <TextInput
              source="normal_ref_string"
              label="Normal Reference String"
              readOnly={props.readonly}
              fullWidth
              helperText="Use this for qualitative reference values like 'negative', 'positive', '1+', etc. Leave empty if using numeric low/high ranges."
              sx={{
                "& .MuiOutlinedInput-root": {
                  borderRadius: 2,
                  transition: "all 0.2s ease",
                  ...(!props.readonly && {
                    "&:hover": {
                      boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                    },
                  }),
                },
              }}
            />

            <Box sx={{ mt: 2 }}>
              <Typography
                variant="body2"
                sx={{
                  fontWeight: 500,
                  color: theme.palette.text.secondary,
                  mb: 2,
                  display: "flex",
                  alignItems: "center",
                  gap: 1,
                }}
              >
                🎯 Specific Reference Ranges (Age/Gender Based)
                <Chip
                  label="Optional"
                  size="small"
                  color="info"
                  variant="outlined"
                  sx={{ fontSize: "0.7rem" }}
                />
              </Typography>
              <Typography
                variant="caption"
                sx={{
                  color: theme.palette.text.secondary,
                  display: "block",
                  mb: 2,
                }}
              >
                Add age and gender-specific reference ranges. If no criteria matches, the default ranges above will be used.
              </Typography>

              <ArrayInput source="specific_ref_ranges">
                <SimpleFormIterator
                  inline={false}
                  disableReordering={props.readonly}
                  disableAdd={props.readonly}
                  disableRemove={props.readonly}
                  sx={{
                    "& .RaSimpleFormIterator-line": {
                      border: `1px solid ${theme.palette.divider}`,
                      borderRadius: 2,
                      p: 2,
                      mb: 2,
                      backgroundColor: theme.palette.background.default,
                    },
                  }}
                >
                  <Stack spacing={2} width="100%">
                    <Stack direction="row" gap={2}>
                      <AutocompleteInput
                        source="gender"
                        label="Gender"
                        choices={[
                          { id: "", name: "All Genders" },
                          { id: "M", name: "Male" },
                          { id: "F", name: "Female" },
                        ]}
                        readOnly={props.readonly}
                        fullWidth
                        sx={{
                          "& .MuiOutlinedInput-root": {
                            borderRadius: 1.5,
                          },
                        }}
                      />
                      <NumberInput
                        source="age_min"
                        label="Min Age (years)"
                        helperText="Leave empty for no minimum"
                        readOnly={props.readonly}
                        fullWidth
                        sx={{
                          "& .MuiOutlinedInput-root": {
                            borderRadius: 1.5,
                          },
                        }}
                      />
                      <NumberInput
                        source="age_max"
                        label="Max Age (years)"
                        helperText="Leave empty for no maximum"
                        readOnly={props.readonly}
                        fullWidth
                        sx={{
                          "& .MuiOutlinedInput-root": {
                            borderRadius: 1.5,
                          },
                        }}
                      />
                    </Stack>
                    <Stack direction="row" gap={2}>
                      <NumberInput
                        source="low_ref_range"
                        label="Low Range"
                        readOnly={props.readonly}
                        fullWidth
                        sx={{
                          "& .MuiOutlinedInput-root": {
                            borderRadius: 1.5,
                          },
                        }}
                      />
                      <NumberInput
                        source="high_ref_range"
                        label="High Range"
                        readOnly={props.readonly}
                        fullWidth
                        sx={{
                          "& .MuiOutlinedInput-root": {
                            borderRadius: 1.5,
                          },
                        }}
                      />
                      <TextInput
                        source="normal_ref_string"
                        label="OR Normal String"
                        helperText="e.g. 'Negative', 'Positive'"
                        readOnly={props.readonly}
                        fullWidth
                        sx={{
                          "& .MuiOutlinedInput-root": {
                            borderRadius: 1.5,
                          },
                        }}
                      />
                    </Stack>
                  </Stack>
                </SimpleFormIterator>
              </ArrayInput>
            </Box>
          </Stack>
        </CardContent>
      </Card>

      <Card
        elevation={0}
        sx={{
          border: `1px solid ${theme.palette.divider}`,
          borderRadius: 2,
        }}
      >
        <CardContent sx={{ p: 3 }}>
          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              gap: 1.5,
              mb: 3,
            }}
          >
            <Typography
              variant="subtitle1"
              sx={{
                fontWeight: 600,
                color: theme.palette.text.primary,
              }}
            >
              📋 Additional Settings
            </Typography>
            <Chip
              label="Required"
              size="small"
              color="error"
              variant="outlined"
              sx={{ ml: "auto", fontSize: "0.75rem" }}
            />
          </Box>

          <Stack spacing={3}>
            <NumberInput
              source="decimal"
              readOnly={props.readonly}
              validate={[required()]}
              fullWidth
              sx={{
                "& .MuiOutlinedInput-root": {
                  borderRadius: 2,
                  transition: "all 0.2s ease",
                  ...(!props.readonly && {
                    "&:hover": {
                      boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                    },
                  }),
                },
              }}
            />
            <Box sx={{ gridColumn: "span 4" }}>
              <BooleanInput
                source="is_calculated_test"
                label="Calc Test"
                disabled={props.readonly}
              />
            </Box>
            <Box>
              <Typography
                variant="body1"
                sx={{
                  fontWeight: 500,
                  color: theme.palette.text.secondary,
                  mb: 2,
                }}
              >
                Specimen Types
              </Typography>
              <ArrayInput source="types">
                <SimpleFormIterator inline>
                  <FeatureList
                    source="type"
                    readOnly={props.readonly}
                    types="specimen-type"
                  >
                    <AutocompleteInput
                      source="type"
                      readOnly={props.readonly}
                      sx={{
                        "& .MuiOutlinedInput-root": {
                          borderRadius: 2,
                          transition: "all 0.2s ease",
                          ...(!props.readonly && {
                            "&:hover": {
                              boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                            },
                          }),
                        },
                      }}
                    />
                  </FeatureList>
                </SimpleFormIterator>
              </ArrayInput>
            </Box>

            {/* Show Reference Ranges Summary in readonly mode */}
            {props.readonly && props.mode === "SHOW" && (
              <SpecificReferenceRangesDisplay />
            )}
          </Stack>
        </CardContent>
      </Card>
    </Stack>
  );
}

function TestTypeForm(props: TestTypeFormProps) {
  return (
    <Box sx={{ p: { xs: 2, sm: 3 } }}>
      <SimpleForm
        disabled={props.readonly}
        toolbar={props.readonly === true ? false : undefined}
        warnWhenUnsavedChanges
        sx={{
          "& .RaSimpleForm-form": {
            backgroundColor: "transparent",
            boxShadow: "none",
            padding: 0,
          },
        }}
      >
        <TestTypeInput {...props} />
      </SimpleForm>
    </Box>
  );
}

export function TestTypeEdit() {
  const theme = useTheme();

  return (
    <Box
      sx={{
        minHeight: "100vh",
        bgcolor: theme.palette.background.default,
        pb: 4,
      }}
    >
      <Edit
        mutationMode="pessimistic"
        title="Edit Test Type"
        redirect={"list"}
        transform={(data: any) => {
          // Ensure device_ids is sent to backend
          console.log('Edit transform - sending data:', data);
          return data;
        }}
        mutationOptions={{
          onSuccess: (data: any) => {
            console.log('Mutation success:', data);
          }
        }}
        queryOptions={{
          onSuccess: (data: any) => {
            console.log('Query success - received data:', data);
            // Transform devices array to device_ids if not present
            if (data.devices && Array.isArray(data.devices) && !data.device_ids) {
              data.device_ids = data.devices.map((d: any) => d.id);
              console.log('Transformed device_ids:', data.device_ids);
            }
          }
        }}
      >
        <TestTypeForm readonly={false} mode={"EDIT"} />
        <Box sx={{ px: { xs: 4, sm: 4.5 } }}>
          <ReferenceSection />
        </Box>
      </Edit>
    </Box>
  );
}

export function TestTypeCreate() {
  const theme = useTheme();
  return (
    <Box
      sx={{
        minHeight: "100vh",
        bgcolor: theme.palette.background.default,
        pb: 4,
      }}
    >
      <Create title="Create Test Type" redirect={"list"}>
        <TestTypeForm readonly={false} mode={"CREATE"} />
        <Box sx={{ px: { xs: 4, sm: 4.5 } }}>
          <ReferenceSection />
        </Box>
      </Create>
    </Box>
  );
}

function TestTypeListActions() {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const axios = useAxios();

  const handleUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setUploading(true);
    const formData = new FormData();
    formData.append("file", file);

    try {
      await axios.post("/test-type/upload", formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
    } finally {
      setUploading(false);
      // Reset file input
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <TopToolbar>
      <Button
        label={uploading ? "Uploading..." : "Upload File"}
        onClick={handleUploadClick}
      >
        <CreateNewFolderIcon />
        <input
          type="file"
          ref={fileInputRef}
          onChange={handleUpload}
          style={{ display: "none" }}
          accept=".csv"
        />
      </Button>
      <CreateButton />
    </TopToolbar>
  );
}
