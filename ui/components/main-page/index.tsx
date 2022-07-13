import React, { useEffect, useState } from "react";
import { Toggle, Node } from "../";
import styles from "./styles.module.scss";

const env = process.env.NODE_ENV;

const apiUrl = env == "production" ? "/info" : "http://35.246.94.55/info";

export function MainPage() {
  const [isTypeOfResourceToggled, setIsTypeOfResourceToggled] = useState(false);
  const [isAllocated, setIsAllocated] = useState(false);
  const [hoveredWorkload, setHoveredWorkload] = useState(undefined);
  const [soundOn, setSoundOn] = useState(false);
  const [sound, setSound] = useState<any>(undefined);
  const [response, setResponse] = useState<any>(undefined);

  useEffect(() => {
    if (typeof window !== "undefined") {
      //Set the sound object
      const sound = new Audio("/tetris.mp3");
      sound.volume = 0.1;
      sound.loop = true;
      setSound(sound);
    }
    getData();
  }, []);

  async function getData() {
    //get the data from the API and set it to the state
    try {
      const response = await fetch(apiUrl);
      const json = await response.json();
      setResponse(json);
    } catch (error) {
      console.log(error);
    }
  }

  function toggleSound() {
    //Toggle the sound and the state on/off
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
        onClick={() => toggleSound()}
        className={`${styles.soundIcon} ${soundOn ? styles.soundOn : ""}`}
      ></div>
      <div className={styles.content}>
        <p>
          binpak is a tool that allows kubernetes newbies to visualise core
          concepts like nodes, workloads, requests and limits without the 🤯
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
