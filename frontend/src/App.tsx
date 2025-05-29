import { useState } from "react";
import { SelectFile } from "../bindings/github.com/eryalito/http-file-share/internal/services/dialogsservice";
import { SetFileToServe, GetAddresses } from "../bindings/github.com/eryalito/http-file-share/internal/services/httpfileserverservice";
import './index.css'
import { QRCodeSVG } from "qrcode.react";


function App() {
  const [selectedFile, setSelectedFile] = useState<string | null>(null);
  const [isSharing, setIsSharing] = useState(false);
  const [addresses, setAddresses] = useState<string[]>([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [modalAddress, setModalAddress] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  // Fetch addresses on component mount
  useState(() => {
    GetAddresses().then((response) => {
      console.log("Fetched addresses:", response);
      setAddresses(response);
    }).catch((error) => {
      console.error("Error fetching addresses:", error);
    });
  });


  const openFile = () => {
    SelectFile("File to Share").then((response) => {
      console.log("Selected file:", response);
      if (response) {
        setSelectedFile(response.Path);
      }
    }).catch((error) => {
      console.error("Error selecting file:", error);
      setSelectedFile(null);
    });
  };

  const getFileName = (filePath: string) => {
    const parts = filePath.split(/[/\\]/);
    return parts.length > 0 ? parts[parts.length - 1] : filePath;
  };

  const shareFile = async () => {
    if (!selectedFile) {
      console.error("No file selected to share.");
      return;
    }
    try {
      if (!isSharing) {
        await SetFileToServe(selectedFile);
      } else {
        await SetFileToServe("");
      }
      setIsSharing(!isSharing);
    } catch (error) {
      console.error("Error setting file to serve:", error);
      return;
    }
  };

  const handleAddressClick = (address: string) => {
    setModalAddress(address);
    setModalOpen(true);
  };

  const handleCopy = () => {
    if (modalAddress) {
      navigator.clipboard.writeText(`http://${modalAddress}/file`);
      setCopied(true);
      setTimeout(() => setCopied(false), 500); // Animation lasts 0.5s
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-700 flex flex-col items-center justify-center p-4">
      {/* Modal */}
      {modalOpen && modalAddress && (
        <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
          <div className="bg-slate-800 p-6 rounded-xl shadow-2xl w-full max-w-md relative">
            <button
              className="absolute top-2 right-2 text-slate-400 hover:text-slate-200 text-xl"
              onClick={() => setModalOpen(false)}
              aria-label="Close"
            >
              &times;
            </button>
            <h2 className="text-lg font-semibold text-cyan-300 mb-4">Share URL</h2>
            <div className="mb-4">
              <span className="font-mono text-cyan-200 break-all">{`http://${modalAddress}/file`}</span>
            </div>
            <div className="flex justify-center mb-4">
              <QRCodeSVG
                value={`http://${modalAddress}/file`}
                size={180}
                bgColor="#1e293b"
                fgColor="#67e8f9"
                className="rounded-lg shadow-lg"
              />
            </div>
            <div className="flex gap-3">
              <button
                onClick={handleCopy}
                className={`!w-full bg-cyan-600 hover:bg-cyan-700 text-white px-4 py-2 rounded-lg font-semibold transition 
                  ${copied ? "animate-pulse scale-105" : ""}`}
              >
                {copied ? "Copied!" : "Copy URL"}
              </button>
            </div>
          </div>
        </div>
      )}
      <div className="bg-white/10 backdrop-blur-md shadow-2xl rounded-xl p-8 md:p-12 w-full max-w-lg transform transition-all duration-500">
        <h1 className="text-4xl md:text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-600 mb-8 text-center">
          Share Your Files
        </h1>
        {!selectedFile && (
          <div className="flex flex-col items-center justify-center">
            <button
              onClick={openFile}
              className="!w-full bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 text-white font-semibold py-3 px-6 rounded-lg text-lg shadow-lg hover:shadow-xl transform transition-transform duration-300 ease-in-out hover:scale-105 focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:ring-opacity-75"
            >
              Select a File to Share
            </button>
            <p className="text-slate-300 text-sm mt-4">Click the button above to select a file.</p>
          </div>
        )}
        {selectedFile && (
          <>
            <div className="mt-8 p-6 bg-slate-800/50 rounded-lg shadow-inner animate-fadeIn">
              <p className="text-slate-300 text-sm mb-1">Ready to share:</p>
              <p className="text-lg font-medium text-cyan-400 break-all mb-2">{getFileName(selectedFile)}</p>
              {isSharing && (
                <div className="text-sm text-slate-400 mb-2">
                  <p className="mb-1">Sharing on:</p>
                  <ul className="space-y-3 mt-2">
                    {addresses.map((address, index) => (
                      <li
                        key={index}
                        className="flex items-center bg-slate-900/70 rounded-lg px-4 py-2 shadow hover:shadow-lg transition-shadow duration-200 group"
                      >
                        <span className="mr-3 text-cyan-400 group-hover:text-cyan-300 transition-colors duration-200">
                          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V4a2 2 0 10-4 0v1.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                          </svg>
                        </span>
                        <span
                          onClick={() => handleAddressClick(address)}
                          className="text-cyan-300 group-hover:text-cyan-200 font-mono hover:underline break-all transition-colors duration-200 cursor-pointer"
                        >
                          {address}/file
                        </span>
                        <span className="ml-auto">
                          <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-slate-400 group-hover:text-cyan-300 transition-colors duration-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 5l7 7m0 0l-7 7m7-7H3" />
                          </svg>
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
              <button
                onClick={shareFile}
                className={`!w-full ${isSharing
                  ? "bg-gradient-to-r from-red-500 to-pink-600 hover:from-red-600 hover:to-pink-700 focus:ring-red-400"
                  : "bg-gradient-to-r from-green-500 to-teal-600 hover:from-green-600 hover:to-teal-700 focus:ring-teal-400"
                  } text-white font-semibold py-3 px-6 rounded-lg text-lg shadow-lg hover:shadow-xl transform transition-transform duration-300 ease-in-out hover:scale-105 focus:outline-none focus:ring-2 focus:ring-opacity-75`}
              >
                {isSharing ? "Stop Sharing" : "Share This File"}
              </button>
            </div>
            <div className="mt-4 text-center">
              <p className="text-slate-300 text-sm hover:scale-105 transition-transform duration-300 cursor-pointer underline"
                onClick={() => {setSelectedFile(null); SetFileToServe("");}}>
                Select another file to share
              </p>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default App;
