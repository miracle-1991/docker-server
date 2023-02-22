#!/usr/bin/python3
# -*- coding: UTF-8 -*-

from keplergl import KeplerGl
import pandas as pd
from pandas import DataFrame

defaultconfig = {
    'version': 'v1',
    'config': {
        'visState': {
            'layers': [{
                'type': 'point',
                'visualChannels': {
                    'sizeField': {
                        'type': 'integer',
                        'name': 'value'
                    },
                    'coverageField': None,
                    'colorScale': 'quantize',
                    'coverageScale': 'linear',
                    'colorField': {
                        'type': 'integer',
                        'name': 'value'
                    },
                    'sizeScale': 'linear'
                },
                'config': {
                    'dataId': 'data_1',
                    'color': [250, 116, 0],
                    'textLabel': {
                        'color': [255, 255, 255],
                        'field': {
                            'name': "Point1",
                            'displayName': 'Point1'
                        },
                        'size': 50,
                        'anchor': 'middle',
                        'offset': [0, 0]
                    },
                    'label': 'H3 Hexagon',
                    'isVisible': True,
                    'visConfig': {
                        'coverageRange': [0, 1],
                        'opacity': 0.8,
                        'elevationScale': 5,
                        'hi-precision': False,
                        'coverage': 1,
                        'enable3d': True,
                        'sizeRange': [0, 500],
                        'colorRange': {
                            'category': 'Uber',
                            'type': 'sequential',
                            'colors': ['#194266', '#355C7D', '#63617F', '#916681', '#C06C84', '#D28389', '#E59A8F', '#F8B195'],
                            'reversed': False,
                            'name': 'Sunrise 8'
                        }
                    },
                    'columns': {
                        'hex_id': 'hex_id'
                    }
                },
                'id': 'jdys7lp'
            },{
                'type': 'point',
                'visualChannels': {
                    'sizeField': {
                        'type': 'integer',
                        'name': 'value'
                    },
                    'coverageField': None,
                    'colorScale': 'quantize',
                    'coverageScale': 'linear',
                    'colorField': {
                        'type': 'integer',
                        'name': 'value'
                    },
                    'sizeScale': 'linear'
                },
                'config': {
                    'dataId': 'data_1',
                    'color': [250, 116, 0],
                    'textLabel': {
                        'color': [255, 255, 255],
                        'field': {
                            'name': "Point1",
                            'displayName': 'Point1'
                        },
                        'size': 50,
                        'anchor': 'middle',
                        'offset': [0, 0]
                    },
                    'label': 'H3 Hexagon',
                    'isVisible': True,
                    'visConfig': {
                        'coverageRange': [0, 1],
                        'opacity': 0.8,
                        'elevationScale': 5,
                        'hi-precision': False,
                        'coverage': 1,
                        'enable3d': True,
                        'sizeRange': [0, 500],
                        'colorRange': {
                            'category': 'Uber',
                            'type': 'sequential',
                            'colors': ['#194266', '#355C7D', '#63617F', '#916681', '#C06C84', '#D28389', '#E59A8F', '#F8B195'],
                            'reversed': False,
                            'name': 'Sunrise 8'
                        }
                    },
                    'columns': {
                        'hex_id': 'hex_id'
                    }
                },
                'id': 'jdys7ld'
            }],
            'interactionConfig': {
                'brush': {
                    'enabled': False,
                    'size': 0.5
                },
                'tooltip': {
                    'fieldsToShow': {
                        'data_1': ['hex_id', 'value']
                    },
                    'enabled': True
                }
            },
            'splitMaps': [],
            'layerBlending': 'normal',
            'filters': []
        },
        'mapState': {
            'bearing': 2.6192893401015205,
            'dragRotate': True,
            'zoom': 12.32053899007826,
            'longitude': -122.42590232651203,
            'isSplit': False,
            'pitch': 37.374216241015446,
            'latitude': 37.76209132041332
        },
        'mapStyle': {
            'mapStyles': {},
            'topLayerGroups': {},
            'styleType': 'dark',
            'visibleLayerGroups': {
                'building': True,
                'land': True,
                '3d building': False,
                'label': True,
                'water': True,
                'border': False,
                'road': True
            }
        }
    }
}

class KeplerDraw:
    def __init__(self, filepath, groundtruthcsv, currentcsv, centerLat, centerLon, outpath):
        self.config = defaultconfig
        if groundtruthcsv.endswith("csv") == True:
            df1 = pd.read_csv(filepath + "/" + groundtruthcsv)
            df1name = groundtruthcsv.replace(".csv", "")
        else:
            data = {
                "lat": [centerLat],
                "lon": [centerLon]
            }
            df1 = DataFrame(data)
            df1name = "GroundTruth"
        df2 = pd.read_csv(filepath + "/" + currentcsv)
        df2name = currentcsv.replace(".csv", "")
        self.config["config"]["mapState"]["latitude"]   = centerLat
        self.config["config"]["mapState"]["longitude"]  = centerLon
        self.config["config"]["visState"]["layers"][0]["config"]["dataId"]   = df1name
        self.config["config"]["visState"]["layers"][0]["config"]["textLabel"]["field"]["name"]          = df1name
        self.config["config"]["visState"]["layers"][0]["config"]["textLabel"]["field"]["displayName"]   = df1name
        self.config["config"]["visState"]["layers"][1]["config"]["dataId"]   = df2name
        self.config["config"]["visState"]["layers"][1]["config"]["textLabel"]["field"]["name"]          = df1name
        self.config["config"]["visState"]["layers"][1]["config"]["textLabel"]["field"]["displayName"]   = df1name
        self.map = KeplerGl(
            data={
                df1name: df1,
                df2name: df2
            },
            config=self.config
        )
        self.htmlfile = outpath + "/" + currentcsv.replace(".csv", ".html")

    def WriteToHTML(self):
        self.map.save_to_html(file_name=self.htmlfile)
        return self.htmlfile


if __name__ == '__main__':
    k = KeplerDraw(
        filepath="/Users/xiaolong.ji/Downloads/rtK/20230215/result/forestrouteloop1/Note20Ultra",
        groundtruthcsv="decodedongle-donglereplacement-driver13770990-dongle.csv",
        currentcsv="PMODE_DGPS_GPS_GLO_GAL.csv",
        centerLat=1.2911047,
        centerLon=103.7929342
    )
    k.getHTML()
