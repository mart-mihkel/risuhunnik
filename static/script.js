import * as THREE from "three";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";

const NORD0 = 0x2e3440;
const NORD6 = 0xeceff4;
const LOADER = new GLTFLoader();

// TODO: intersection observer for lazy load

document
  .getElementById("joke-container")
  .addEventListener("htmx:afterSwap", () => {
    initFlounder();
  });

function initFlounder() {
  const div = document.getElementById("joke-lest"); // TODO: error on no div
  const camera = new THREE.PerspectiveCamera();
  camera.position.set(0, 0, 5);

  const renderer = new THREE.WebGLRenderer();
  renderer.setClearColor(NORD0);

  const scene = new THREE.Scene();
  const light = new THREE.AmbientLight(NORD6, 1);
  scene.add(light);

  /** @type {THREE.Object3D | undefined} */ let model;
  LOADER.load(
    "models/flounder.glb",
    (gltf) => {
      model = gltf.scene;
      model.position.x = 2.5;
      scene.add(model);
    },
    (xhr) => {
      // TODO: loading screen
      }
    },
    (error) => {
      // TODO: error handling
      console.error(error);
    },
  );

  renderer.setAnimationLoop(() => {
    if (!model) {
      return;
    }

    // TODO: bezier/cubic spline
    model.position.x = Math.max(0, model.position.x - 0.01);
    model.rotation.y += 0.01;
    renderer.render(scene, camera);
  });

  resize(div, camera, renderer);
  div.addEventListener("resize", () => resize(div, camera, renderer));
  div.appendChild(renderer.domElement);
}

/**
 * @param {HTMLDivElement} div
 * @param {THREE.Camera} cam
 * @param {THREE.WebGLRenderer} rend
 */
function resize(div, cam, rend) {
  cam.aspect = div.clientWidth / div.clientHeight;
  cam.updateProjectionMatrix();
  rend.setSize(div.clientWidth, div.clientHeight);
}
