import * as THREE from "three";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";

const NORD0 = 0x2e3440;
const NORD6 = 0xeceff4;

const PI2 = Math.PI / 2;

// TODO: intersection observer for lazy load

document
  .getElementById("joke-container")
  .addEventListener("htmx:afterSwap", () => {
    init(
      "lest",
      "models/flounder.glb",
      (
        /** @type {THREE.Scene} */ scene,
        /** @type {THREE.Object3D} */ model,
      ) => {
        model.position.x = 4;
        model.rotateX(PI2);
        model.scale.set(2, 2, 2);
        scene.add(model);
      },
      (/** @type {THREE.Object3D | undefined} */ model) => {
        if (model === undefined) {
          return;
        }

        // TODO: bezier/cubic spline
        model.position.x = Math.max(0, model.position.x - 0.01);
        model.rotateZ(0.01);
      },
    );

    init(
      "hernes",
      "models/peapod.glb",
      (
        /** @type {THREE.Scene} */ scene,
        /** @type {THREE.Object3D} */ model,
      ) => {
        model.position.x = 2.5;
        model.scale.set(0.4, 0.4, 0.4);
        scene.add(model);
      },
      (/** @type {THREE.Object3D | undefined} */ model) => {
        if (model === undefined) {
          return;
        }

        // TODO: bezier/cubic spline
        model.position.x = Math.max(0, model.position.x - 0.01);
        model.rotation.y += 0.01;
      },
    );
  });

/**
 * @param {string} id
 * @param {string} modelPath
 * @param {function(THREE.Scene, THREE.Object3D): void} modelInit
 * @param {function(THREE.Object3D | undefined): void} animate
 */
function init(id, modelPath, modelInit, animate) {
  const div = document.getElementById(`joke-${id}-visual`);
  if (div === null) {
    throw Error(`Can't init ${id}, no div for visual`);
  }

  const scene = new THREE.Scene();
  scene.add(new THREE.AmbientLight(NORD6, 2));

  const renderer = new THREE.WebGLRenderer();
  renderer.setClearColor(NORD0);

  const camera = new THREE.PerspectiveCamera();
  camera.position.set(0, 5, 5);
  camera.lookAt(new THREE.Vector3());

  /** @type {THREE.Object3D | undefined} */
  let model;
  new GLTFLoader().load(
    modelPath,
    (gltf) => {
      model = gltf.scene;
      modelInit(scene, gltf.scene);
    },
    (xhr) => modelLoading(div, `joke-${id}-loading`, xhr.loaded / xhr.total),
    (error) => modelError(div, `joke-${id}-loading`, error),
  );

  renderer.setAnimationLoop(() => {
    animate(model);
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

/**
 * @param {HTMLDivElement} parent
 * @param {string} id
 * @param {number} prog
 */
function modelLoading(parent, id, prog) {
  if (prog < 1) {
    return;
  }

  const loading = document.getElementById(id);
  parent.removeChild(loading);
}

/**
 * @param {HTMLDivElement} parent
 * @param {string} id
 * @param {Error} error
 */
function modelError(parent, id, error) {
  const loading = document.getElementById(id);
  parent.removeChild(loading);
  console.error(error);
}
