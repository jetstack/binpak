import React, { useEffect, useState } from "react";
import { Toggle, Node } from "../";
import styles from "./styles.module.scss";

const env = process.env.NODE_ENV;

const apiUrl = env == "production" ? "/api/info" : "http://35.246.94.55/info";

export function MainPage() {
  const [isTypeOfResourceToggled, setIsTypeOfResourceToggled] = useState(false);
  const [isAllocated, setIsAllocated] = useState(false);
  const [hoveredWorkload, setHoveredWorkload] = useState(undefined);
  const [soundOn, setSoundOn] = useState(false);
  const [sound, setSound] = useState<any>(undefined);
  const [response, setResponse] = useState<any>(undefined);

  useEffect(() => {
    if (typeof window !== "undefined") {
      const sound = new Audio("/tetris.mp3");
      sound.volume = 0.1;
      sound.loop = true;
      setSound(sound);
    }
    getData();
  }, []);

  async function getData() {
    try {
      const response = await fetch(apiUrl);
      const json = await response.json();
      setResponse(json);
    } catch (error) {
      console.log(error);
    }
  }

  function toggleSound() {
    if (soundOn) {
      sound.pause();
    } else {
      sound.play();
    }
    setSoundOn(!soundOn);
  }

  const isCPU = isTypeOfResourceToggled;

  return (
    <div className={styles.container}>
      <div className={styles.logo}>
        <img src="./images/binpak.svg" alt="" />
      </div>
      <div
        onClick={() => toggleSound(soundOn)}
        className={`${styles.soundIcon} ${soundOn ? styles.soundOn : ""}`}
      ></div>
      <div className={styles.content}>
        <p>
          This is an introduction text to explain what binpak is. This is an
          introduction text to explain what binpak is. This is an introduction
          text to explain what binpak is. This is an introduction text to
          explain what binpak is.
        </p>
      </div>

      <div className={styles.togglesContainer}>
        <Toggle
          leftSideItem="Memory"
          rightSideItem="CPU"
          isToggled={isTypeOfResourceToggled}
          onToggle={setIsTypeOfResourceToggled}
        />
        <div className={styles.divider}></div>
        <Toggle
          leftSideItem="Requested"
          rightSideItem="Limits"
          isToggled={isAllocated}
          onToggle={setIsAllocated}
        />
      </div>
      {response && (
        <>
          <h3>Cluster name: {response.clusterName}</h3>
          <div className={styles.list}>
            {response.instances.map((instance: any) => (
              <Node
                key={instance.name}
                workloads={instance.workloads}
                limit={
                  isCPU ? instance.capacity.cpuM : instance.capacity.memoryMi
                }
                isCPU={isCPU}
                isAllocated={isAllocated}
                hoveredWorkload={hoveredWorkload}
                setHoveredWorkload={setHoveredWorkload}
              />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
