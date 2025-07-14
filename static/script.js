import * as THREE from "three";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";

const NORD0 = 0x2e3440;
const NORD6 = 0xeceff4;

const PI2 = Math.PI / 2;
const INIT_STATE = {};

const jokeContainer = document.getElementById("joke-container");
const observer = new IntersectionObserver(
  (entries) => {
    entries
      .filter((e) => e.isIntersecting)
      .forEach((e) => {
        switch (e.target.id) {
          case "joke-lest-visual":
            initLest();
            break;
          case "joke-hernes-visual":
            initHernes();
            break;
          // TODO: other visuals
        }
      });
  },
  {
    threshold: 0.25,
  },
);

jokeContainer.addEventListener("htmx:afterSwap", () => {
  const visuals = jokeContainer.getElementsByClassName("joke-visual");
  for (let i = 0; i < visuals.length; i++) {
    const child = visuals.item(i);
    observer.observe(child);
  }
});

function initLest() {
  init(
    "lest",
    "models/flounder.glb",
    (/** @type {THREE.Scene} */ scene, /** @type {THREE.Object3D} */ model) => {
      model.position.set(100, 0, 0);
      model.scale.set(2, 2, 2);
      model.rotateX(PI2);
      scene.add(model);
    },
    (
      /** @type {THREE.Clock} */ clock,
      /** @type {THREE.Object3D | undefined} */ model,
    ) => {
      if (model === undefined) {
        return;
      }

      const t = clock.getElapsedTime();
      if (t < 4) {
        model.position.x = 4 - 4 * bezier(t * 0.25);
      }

      model.rotateZ(0.01);
    },
  );
}

function initHernes() {
  init(
    "hernes",
    "models/peapod.glb",
    (/** @type {THREE.Scene} */ scene, /** @type {THREE.Object3D} */ model) => {
      model.scale.set(0, 0, 0);
      scene.add(model);
    },
    (
      /** @type {THREE.Clock} */ clock,
      /** @type {THREE.Object3D | undefined} */ model,
    ) => {
      if (model === undefined) {
        return;
      }

      const t = clock.getElapsedTime();
      if (t < 4) {
        const scale = 0.3 * bezier(t * 0.25);
        model.scale.set(scale, scale, scale);
      }

      model.rotateY(0.01);
    },
  );
}

/**
 * @param {string} id
 * @param {string} modelPath
 * @param {function(THREE.Scene, THREE.Object3D): void} modelInit
 * @param {function(THREE.Object3D | undefined): void} animate
 */
function init(id, modelPath, modelInit, animate) {
  if (id in INIT_STATE) {
    return;
  }

  const div = document.getElementById(`joke-${id}-visual`);
  if (div === null) {
    throw Error(`Can't init ${id}, no div for visual`);
  }

  const clock = new THREE.Clock();
  const scene = new THREE.Scene();
  const light = new THREE.DirectionalLight(NORD6, 2);
  light.position.set(-5, 5, 5);
  scene.add(light);

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
    animate(clock, model);
    renderer.render(scene, camera);
  });

  resize(div, camera, renderer);
  div.addEventListener("resize", () => resize(div, camera, renderer));
  div.appendChild(renderer.domElement);

  INIT_STATE[id] = true;
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

/**
 * ease in out [0, 1] -> [0, 1]
 * @param {number} t
 */
function bezier(t) {
  return t * t * (3.0 - 2.0 * t);
}
