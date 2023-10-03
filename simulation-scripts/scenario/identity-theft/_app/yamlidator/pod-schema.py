{
    "apiVersion": {
        "required": True,
        "type": "string",
        "regex": "v1"
    },
    "kind": {
        "required": True,
        "type": "string",
        "regex": "Pod"
    },
    "metadata": {
        "type":"dict",
        "schema": {
            "name": {
                "type": "string",
            },
            "namespace": {
                "type": "string",
            },
            "labels": {
                "type": "dict",
            },
        }
    },
    "spec": {
        "required": True,
        "type":"dict",
        "schema": {
            "containers": {
                "required": True,
                "type":"list",
                "items": [{
                    "type":"dict",
                    "schema": {
                        "image": {
                            "required": True,
                            "type": "string"
                        },
                        "command": {
                            "type":"list",
                        },
                        "name": {
                            "required": True,
                            "type": "string"
                        },
                        "imagePullPolicy": {
                            "allowed": [
                                "IfNotPresent"
                            ]
                        },
                        "ports": {
                            "type":"list",
                            "items": [{
                                "type":"dict",
                                "schema": {
                                    "name": {
                                        "type": "string"
                                    },
                                    "containerPort": {
                                        "type": "integer"
                                    }
                                }
                            }]
                        },
                       "securityContext": {
                            "type":"dict",
                            "schema": {
                                "allowPrivilegeEscalation": {
                                    "type": "boolean"
                                },
                                "capabilities": {
                                    "type":"dict",
                                    "schema": {
                                        "drop": {
                                        "type":"list",
                                        "items": [{
                                            "allowed": [
                                                "ALL"
                                            ]
                                        }]
                                        },
                                        "add": {
                                        "type":"list",
                                        "items": [{
                                            "allowed": [
                                                "ALL"
                                            ]
                                        }]
                                        }
                                    }
                                },
                                "runAsNonRoot": {
                                    "type": "boolean"
                                },
                                "runAsUser": {
                                    "type": "integer"
                                },
                                "runAsGroup": {
                                    "type": "integer"
                                },
                                "fsGroup": {
                                    "type": "integer"
                                },
                            }
                        },
                       "volumeMounts": {
                            "type":"list",
                            "minlength": 1,
                            "maxlength": 10,
                            "schema": {
                                "type":"dict",
                                "schema": {
                                    "readOnly": {
                                        "type": "boolean"
                                    },
                                    "mountPath": {
                                        "type": "string"
                                    },
                                    "subPath": {
                                        "type": "string"
                                    },
                                    "name": {
                                        "type": "string"
                                    },
                                }
                            }
                       },
                       "env": {
                            "type":"list",
                            "minlength": 1,
                            "maxlength": 10,
                            "schema": {
                                "type":"dict",
                                "schema": {
                                    "name": {
                                        "type": "string"
                                    },
                                    "value": {
                                        "type": "string"
                                    },
                                    "valueFrom": {
                                        "type":"dict",
                                        "schema": {
                                            "secretKeyRef": {
                                                "type":"dict",
                                                "schema": {
                                                    "name": {
                                                        "type": "string"
                                                    },
                                                    "key": {
                                                        "type": "string"
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        },
                       "resources": {
                            "type":"dict",
                            "schema": {
                                "limits": {
                                    "type":"dict",
                                    "schema": {
                                        "memory": {
                                            "type": "string"
                                        },
                                        "cpu": {
                                            "type": "string"
                                        }
                                    }
                                },
                                "requests": {
                                    "type":"dict",
                                    "schema": {
                                        "memory": {
                                            "type": "string"
                                        },
                                        "cpu": {
                                            "type": "string"
                                        }
                                    }
                                }
                            }
                        },
                    }
                }]
            },
            "restartPolicy": {
                "allowed": [
                    "Always"
                ]
            },
            "serviceAccountName": {
                "type": "string"
            },
            "automountServiceAccountToken": {
                "type": "boolean"
            },
            "nodeName": {
                "type": "string"
            },
            "hostname": {
                "type": "string"
            },
            "nodeSelector": {
                "type":"dict",
                "schema": {
                    "additionalProperties": {
                        "type": "string"
                    }
                }
            },
            "hostIPC": {
                "type": "boolean"
            },
            "hostPID": {
                "type": "boolean"
            },
            "terminationGracePeriodSeconds": {
                "type": "integer"
            },
            "dnsPolicy": {
                "type": "string"
            },
            "initContainers": {
                "type":"list",
                "items": [{
                    "type":"dict",
                    "schema": {
                        "image": {
                            "required": True,
                            "type": "string"
                        },
                        "command": {
                            "type":"list",
                        },
                        "name": {
                            "required": True,
                            "type": "string"
                        },
                        "imagePullPolicy": {
                            "allowed": [
                                "IfNotPresent"
                            ]
                        },
                        "ports": {
                            "type":"list",
                            "items": [{
                                "type":"dict",
                                "schema": {
                                    "name": {
                                        "type": "string"
                                    },
                                    "containerPort": {
                                        "type": "integer"
                                    }
                                }
                            }]
                        },
                       "securityContext": {
                            "type":"dict",
                            "schema": {
                                "allowPrivilegeEscalation": {
                                    "type": "boolean"
                                },
                                "capabilities": {
                                    "type":"dict",
                                    "schema": {
                                        "drop": {
                                        "type":"list",
                                        "items": [{
                                            "allowed": [
                                                "ALL"
                                            ]
                                        }]
                                        },
                                        "add": {
                                        "type":"list",
                                        "items": [{
                                            "allowed": [
                                                "ALL"
                                            ]
                                        }]
                                        }
                                    }
                                },
                                "runAsNonRoot": {
                                    "type": "boolean"
                                },
                                "runAsUser": {
                                    "type": "integer"
                                },
                                "runAsGroup": {
                                    "type": "integer"
                                },
                                "fsGroup": {
                                    "type": "integer"
                                },
                            }
                        },
                       "volumeMounts": {
                            "type":"list",
                            "minlength": 1,
                            "maxlength": 10,
                            "schema": {
                                "type":"dict",
                                "schema": {
                                    "readOnly": {
                                        "type": "boolean"
                                    },
                                    "mountPath": {
                                        "type": "string"
                                    },
                                    "subPath": {
                                        "type": "string"
                                    },
                                    "name": {
                                        "type": "string"
                                    },
                                }
                            }
                       },
                       "env": {
                            "type":"list",
                            "minlength": 1,
                            "maxlength": 10,
                            "schema": {
                                "type":"dict",
                                "schema": {
                                    "name": {
                                        "type": "string"
                                    },
                                    "value": {
                                        "type": "string"
                                    },
                                    "valueFrom": {
                                        "type":"dict",
                                        "schema": {
                                            "secretKeyRef": {
                                                "type":"dict",
                                                "schema": {
                                                    "name": {
                                                        "type": "string"
                                                    },
                                                    "key": {
                                                        "type": "string"
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        },
                       "resources": {
                            "type":"dict",
                            "schema": {
                                "limits": {
                                    "type":"dict",
                                    "schema": {
                                        "memory": {
                                            "type": "string"
                                        },
                                        "cpu": {
                                            "type": "string"
                                        }
                                    }
                                },
                                "requests": {
                                    "type":"dict",
                                    "schema": {
                                        "memory": {
                                            "type": "string"
                                        },
                                        "cpu": {
                                            "type": "string"
                                        }
                                    }
                                }
                            }
                        },
                    }
                }]
            },
            "imagePullSecrets": {
                "type": "dict",
                "schema": {
                    "name": {
                        "type": "string"
                    }
                }
            },
            "volumes": {
                "type":"list",
                "items": [{
                    "type":"dict",
                    "schema": {
                        "name": {
                            "type": "string"
                        },
                        "emptyDir": {
                            "type":"list",
                            "items": [{
                                "type":"dict",
                                "schema": {
                                    "sizeLimit": {
                                        "type": "string"
                                    },
                                    "medium": {
                                        "allowed": [
                                            "Memory"
                                        ]
                                    },
                                }
                            }]
                        },
                        "hostPath": {
                                "type":"dict",
                                "schema": {
                                    "path": {
                                        "type": "string"
                                    },
                                    "type": {
                                        "allowed": [
                                            "Directory",
                                            "DirectoryOrCreate",
                                            "FileOrCreate"
                                        ]
                                    }
                                }
                        },
                        "persistentVolumeClaim": {
                            "type":"list",
                            "items": [{
                                "type":"dict",
                                "schema": {
                                    "claimName": {
                                        "type": "string"
                                    }
                                }
                            }]
                        },
                        "nfs": {
                            "type":"list",
                            "items": [{
                                "type":"dict",
                                "schema": {
                                    "server": {
                                        "type": "string"
                                    },
                                    "path": {
                                        "type": "string"
                                    },
                                    "readOnly": {
                                        "type": "boolean"
                                    }
                                }
                            }]
                        }
                    }
                }]
            }
        }
    }
}