mapboxgl.accessToken = 'pk.eyJ1IjoibWVuZGV6NjQxMiIsImEiOiJjanI4aXFsd3IwN2VkNDRwZG9sanUybWZkIn0.ZoqHstzqGbjETTzoDqWjeg';

var map = new mapboxgl.Map({
    container: 'heatmap',
    center: [-78.8986, 35.9940],
    zoom: 13,
    style: 'mapbox://styles/mapbox/streets-v9',
  });

map.on('render', (e) => {
    // TODO: Add debounce time here so it's not constantly updating boundaries
    console.log('event: ', e);
    const boundaries = e.target.getBounds();
    console.log('bounds: ', boundaries);
})