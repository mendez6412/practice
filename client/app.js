mapboxgl.accessToken = 'pk.eyJ1IjoibWVuZGV6NjQxMiIsImEiOiJjanI4aXFsd3IwN2VkNDRwZG9sanUybWZkIn0.ZoqHstzqGbjETTzoDqWjeg';
var sourceCounter = 0;
var lastBounds = { _sw: { lat: null, ng: null}, _ne: { lat: null, lng: null } };

var map = new mapboxgl.Map({
    container: 'heatmap',
    center: [-78.8986, 35.9940],
    zoom: 13,
    style: 'mapbox://styles/mapbox/streets-v9',
  });

function hasChangedBounds(bounds, lastBounds) {
    return bounds._sw.lat === lastBounds._sw.lat && bounds._sw.lng === lastBounds._sw.lng && bounds._ne.lat === lastBounds._ne.lat && bounds._ne.lng === lastBounds._ne.lng;
}

function getAddressesByBoundary(e) {
    const bounds = e.target.getBounds();
    if (!hasChangedBounds(bounds, lastBounds)) {
        lastBounds = bounds;
        var sourceId = 'addresses' + sourceCounter;
        if (sourceCounter > 0) {
            map.removeLayer('addresses-heat');
            map.removeLayer('addresses-point');
        }
        const url = `http://localhost:8000/getAddressesByBoundary/${bounds._sw.lat}/${bounds._sw.lng}/${bounds._ne.lat}/${bounds._ne.lng}`;
        fetch(url)
        .then((response) => {
            return response.json();
        })
        .then((jsonResponse) => {
            addHeatmap(sourceId, jsonResponse);
            sourceCounter++;
        })
    }
}

map.on('render', _.throttle((e) => {getAddressesByBoundary(e)}, 2500));

function addHeatmap(sourceId, data) {
    map.addSource(sourceId, {
        type: 'geojson',
        data
    });

    map.addLayer({
        id: 'addresses-heat',
        type: 'heatmap',
        source: sourceId,
        maxzoom: 15,
        paint: {
            // increase weight as diameter breast height increases
            'heatmap-weight': {
              property: 'dbh',
              type: 'exponential',
              stops: [
                [1, 0],
                [62, 1]
              ]
            },
            // increase intensity as zoom level increases
            'heatmap-intensity': {
              stops: [
                [11, 1],
                [15, 3]
              ]
            },
            // assign color values be applied to points depending on their density
            'heatmap-color': [
              'interpolate',
              ['linear'],
              ['heatmap-density'],
              0, 'rgba(236,222,239,0)',
              0.2, 'rgb(208,209,230)',
              0.4, 'rgb(166,189,219)',
              0.6, 'rgb(103,169,207)',
              0.8, 'rgb(28,144,153)'
            ],
            // increase radius as zoom increases
            'heatmap-radius': {
              stops: [
                [11, 15],
                [15, 20]
              ]
            },
            // decrease opacity to transition into the circle layer
            'heatmap-opacity': {
              default: 1,
              stops: [
                [14, 1],
                [15, 0]
              ]
            },
        }
    }, 'waterway-label');

    map.addLayer({
        id: 'addresses-point',
        type: 'circle',
        source: sourceId,
        minzoom: 14,
        paint: {
            // increase the radius of the circle as the zoom level and dbh value increases
            'circle-radius': {
                property: 'dbh',
                type: 'exponential',
                stops: [
                [{ zoom: 15, value: 1 }, 5],
                [{ zoom: 15, value: 62 }, 10],
                [{ zoom: 22, value: 1 }, 20],
                [{ zoom: 22, value: 62 }, 50],
                ]
            },
            'circle-color': {
                property: 'dbh',
                type: 'exponential',
                stops: [
                [0, 'rgba(236,222,239,0)'],
                [10, 'rgb(236,222,239)'],
                [20, 'rgb(208,209,230)'],
                [30, 'rgb(166,189,219)'],
                [40, 'rgb(103,169,207)'],
                [50, 'rgb(28,144,153)'],
                [60, 'rgb(1,108,89)']
                ]
            },
            'circle-stroke-color': 'white',
            'circle-stroke-width': 1,
            'circle-opacity': {
                stops: [
                [14, 0],
                [15, 1]
                ]
            }
        }
    }, 'waterway-label');
      
}